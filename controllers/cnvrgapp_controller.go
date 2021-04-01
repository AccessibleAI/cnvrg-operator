package controllers

import (
	"context"
	"fmt"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/controlplane"
	"github.com/cnvrg-operator/pkg/dbs"
	"github.com/cnvrg-operator/pkg/desired"
	"github.com/cnvrg-operator/pkg/logging"
	"github.com/cnvrg-operator/pkg/monitoring"
	"github.com/cnvrg-operator/pkg/networking"
	"github.com/cnvrg-operator/pkg/registry"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
	"io/ioutil"
	v1apps "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"path/filepath"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
	"time"
)

const CnvrgappFinalizer = "cnvrgapp.mlops.cnvrg.io/finalizer"

type CnvrgAppReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var cnvrgAppLog logr.Logger

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *CnvrgAppReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	cnvrgAppLog = r.Log.WithValues("name", req.NamespacedName)
	cnvrgAppLog.Info("starting cnvrgapp reconciliation")

	equal, err := r.syncCnvrgAppSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !equal {
		return ctrl.Result{Requeue: true}, nil
	}

	cnvrgApp, err := r.getCnvrgAppSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cnvrgApp == nil {
		return ctrl.Result{}, nil // probably spec was deleted, no need to reconcile
	}

	// Setup finalizer
	if cnvrgApp.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			cnvrgApp.ObjectMeta.Finalizers = append(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer)
			if err := r.Update(context.Background(), cnvrgApp); err != nil {
				cnvrgAppLog.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer) {
			r.updateStatusMessage(mlopsv1.StatusRemoving, "removing cnvrg spec", cnvrgApp)
			if err := r.cleanup(cnvrgApp); err != nil {
				return ctrl.Result{}, err
			}
			cnvrgApp.ObjectMeta.Finalizers = removeString(cnvrgApp.ObjectMeta.Finalizers, CnvrgappFinalizer)
			if err := r.Update(context.Background(), cnvrgApp); err != nil {
				cnvrgAppLog.Info("error in removing finalizer, checking if cnvrgApp object still exists")
				// if update was failed, make sure that cnvrgApp still exists
				spec, e := r.getCnvrgAppSpec(req.NamespacedName)
				if spec == nil && e == nil {
					return ctrl.Result{}, nil // probably spec was deleted, stop reconcile
				}
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	ready, percentageReady, err := r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}
	if percentageReady == 100 {
		percentageReady = 99
	}
	r.updateStatusMessage(mlopsv1.StatusReconciling, fmt.Sprintf("reconciling... (%d%%)", percentageReady), cnvrgApp)

	if err := r.applyManifests(cnvrgApp); err != nil {
		return ctrl.Result{}, err
	}

	// get control plan readiness
	ready, percentageReady, err = r.getControlPlaneReadinessStatus(cnvrgApp)
	if err != nil {
		return ctrl.Result{}, err
	}

	statusMsg := fmt.Sprintf("successfully reconciled, ready (%d%%)", percentageReady)
	cnvrgAppLog.Info(statusMsg)

	if ready {
		r.updateStatusMessage(mlopsv1.StatusReady, statusMsg, cnvrgApp)
		if err := r.triggerInfraReconciler(cnvrgApp, "add"); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	} else {
		requeueAfter, err := time.ParseDuration("30s")
		if err != nil {
			cnvrgAppLog.Error(err, "wrong duration for requeueAfter")
			return ctrl.Result{}, err
		}
		cnvrgAppLog.Info("stack not ready yet, requeuing...")
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
}

func (r *CnvrgAppReconciler) getControlPlaneReadinessStatus(cnvrgApp *mlopsv1.CnvrgApp) (bool, int, error) {

	readyState := make(map[string]bool)

	// check webapp status
	if cnvrgApp.Spec.ControlPlane.WebApp.Enabled == "true" {
		name := types.NamespacedName{Name: cnvrgApp.Spec.ControlPlane.WebApp.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["webApp"] = ready
	}

	// check sidekiq status
	if cnvrgApp.Spec.ControlPlane.Sidekiq.Enabled == "true" {
		name := types.NamespacedName{Name: "sidekiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["sidekiq"] = ready
	}

	// check searchkiq status
	if cnvrgApp.Spec.ControlPlane.Searchkiq.Enabled == "true" {
		name := types.NamespacedName{Name: "searchkiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["searchkiq"] = ready
	}

	// check systemkiq status
	if cnvrgApp.Spec.ControlPlane.Systemkiq.Enabled == "true" {
		name := types.NamespacedName{Name: "systemkiq", Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["searchkiq"] = ready
	}

	// check postgres status
	if cnvrgApp.Spec.Dbs.Pg.Enabled == "true" {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Pg.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["pg"] = ready
	}

	// check minio status
	if cnvrgApp.Spec.Dbs.Minio.Enabled == "true" {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Minio.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["minio"] = ready
	}

	// check redis status
	if cnvrgApp.Spec.Dbs.Redis.Enabled == "true" {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Redis.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["redis"] = ready
	}

	// check es status
	if cnvrgApp.Spec.Dbs.Es.Enabled == "true" {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Dbs.Es.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckStatefulSetReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["es"] = ready
	}

	// check kibana status
	if cnvrgApp.Spec.Logging.Enabled == "true" && cnvrgApp.Spec.Logging.Kibana.Enabled == "true" {
		name := types.NamespacedName{Name: cnvrgApp.Spec.Logging.Kibana.SvcName, Namespace: cnvrgApp.Namespace}
		ready, err := r.CheckDeploymentReadiness(name)
		if err != nil {
			return false, 0, err
		}
		readyState["kibana"] = ready
	}

	percentageReady := 0

	readyCount := 0

	for _, ready := range readyState {
		if ready {
			readyCount++
		}
	}

	if len(readyState) > 0 {
		percentageReady = readyCount * 100 / len(readyState)
	}

	return readyCount == len(readyState), percentageReady, nil
}

func (r *CnvrgAppReconciler) applyManifests(cnvrgApp *mlopsv1.CnvrgApp) error {

	// registry
	cnvrgInfraLog.Info("applying registry")
	if err := desired.Apply(registry.State(), cnvrgApp, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// dbs
	cnvrgAppLog.Info("applying dbs")
	if err := desired.Apply(dbs.AppDbsState(cnvrgApp), cnvrgApp, r.Client, r.Scheme, cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// controlplane
	cnvrgAppLog.Info("applying controlplane")
	if err := desired.Apply(controlplane.State(cnvrgApp), cnvrgApp, r.Client, r.Scheme, cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// networking
	cnvrgAppLog.Info("applying networking")
	if err := desired.Apply(networking.CnvrgAppNetworkingState(cnvrgApp), cnvrgApp, r.Client, r.Scheme, cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// logging
	cnvrgAppLog.Info("applying logging")
	if err := desired.Apply(logging.CnvrgAppLoggingState(cnvrgApp), cnvrgApp, r.Client, r.Scheme, cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	// grafana dashboards
	cnvrgAppLog.Info("applying grafana dashboards ")
	if err := r.createGrafanaDashboards(cnvrgApp); err != nil {
		return err
	}

	// monitoring
	cnvrgAppLog.Info("applying monitoring")
	if err := desired.Apply(monitoring.AppMonitoringState(cnvrgApp), cnvrgApp, r.Client, r.Scheme, cnvrgAppLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgApp)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) createGrafanaDashboards(cnvrgApp *mlopsv1.CnvrgApp) error {

	if cnvrgApp.Spec.Monitoring.Enabled != "true" {
		cnvrgInfraLog.Info("monitoring disabled, skipping grafana deployment")
		return nil
	}

	basePath := "/pkg/monitoring/tmpl/grafana/dashboards-data/"

	for _, dashboard := range desired.GrafanaAppDashboards {
		if dashboard == "node-exporter.json" {
			fmt.Println("as")
		}
		f, err := pkger.Open(basePath + dashboard)
		if err != nil {
			cnvrgAppLog.Error(err, "error reading path", "path", dashboard)
			return err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			cnvrgAppLog.Error(err, "error reading", "file", dashboard)
			return err
		}
		cm := &v1core.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strings.TrimSuffix(filepath.Base(f.Name()), filepath.Ext(f.Name())),
				Namespace: cnvrgApp.Namespace,
			},
			Data: map[string]string{filepath.Base(f.Name()): string(b)},
		}
		if err := ctrl.SetControllerReference(cnvrgApp, cm, r.Scheme); err != nil {
			cnvrgAppLog.Error(err, "error setting controller reference", "file", f.Name())
			return err
		}
		if err := r.Create(context.Background(), cm); err != nil && errors.IsAlreadyExists(err) {
			cnvrgAppLog.V(1).Info("grafana dashboard already exists", "file", dashboard)
			continue
		} else if err != nil {
			cnvrgAppLog.Error(err, "error reading", "file", dashboard)
			return err
		}
	}

	return nil

}

func (r *CnvrgAppReconciler) triggerInfraReconciler(cnvrgApp *mlopsv1.CnvrgApp, op string) error {

	cnvrgAppInfra := &mlopsv1.CnvrgInfraList{}

	if err := r.List(context.Background(), cnvrgAppInfra); err != nil {
		cnvrgAppLog.Error(err, "can't list CnvrgInfra objects")
		return err
	}

	if len(cnvrgAppInfra.Items) == 0 {
		cnvrgAppLog.Info("no CnvrgInfra objects was deployed, skipping infra reconciler")
		return nil
	}

	name := types.NamespacedName{
		Name:      cnvrgAppInfra.Items[0].Spec.InfraReconcilerCm,
		Namespace: cnvrgAppInfra.Items[0].Spec.InfraNamespace,
	}

	cm := &v1core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
	}

	if err := r.Get(context.Background(), name, cm); err != nil && errors.IsNotFound(err) {
		cnvrgAppLog.Info("infra reconciler cm does not exists, skipping", name, name)
		return nil
	} else if err != nil {
		cnvrgAppLog.Error(err, "can't get cm", "name", name)
		return err
	}

	if op == "add" {
		if cm.Data == nil {
			cm.Data = map[string]string{cnvrgApp.Namespace: cnvrgApp.Name}
		} else {
			cm.Data[cnvrgApp.Namespace] = cnvrgApp.Name
		}
	}
	if op == "remove" {
		delete(cm.Data, cnvrgApp.Namespace)
	}
	if err := r.Update(context.Background(), cm); err != nil {
		cnvrgAppLog.Error(err, "can't update cm", "cm", name)
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) updateStatusMessage(status mlopsv1.OperatorStatus, message string, cnvrgApp *mlopsv1.CnvrgApp) {
	if cnvrgApp.Status.Status == mlopsv1.StatusRemoving {
		cnvrgAppLog.Info("skipping status update, current cnvrg spec under removing status...")
		return
	}
	ctx := context.Background()
	cnvrgApp.Status.Status = status
	cnvrgApp.Status.Message = message
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err := r.Status().Update(ctx, cnvrgApp)
		return err
	})
	if err != nil {
		cnvrgAppLog.Error(err, "can't update status")
	}

	//// This check is to make sure that the status is indeed updated
	//// short reconciliations loop might cause status to be applied but not yet saved into BD
	//// and leads to error: "the object has been modified; please apply your changes to the latest version and try again"
	//// to avoid this error, fetch the object and compare the status
	//statusCheckAttempts := 3
	//for {
	//	cnvrgApp, err := r.getCnvrgAppSpec(types.NamespacedName{Namespace: cnvrgApp.Namespace, Name: cnvrgApp.Name})
	//	if err != nil {
	//		cnvrgAppLog.Error(err, "can't validate status update")
	//	}
	//	cnvrgAppLog.V(1).Info("expected status", "status", status, "message", message)
	//	cnvrgAppLog.V(1).Info("current status", "status", cnvrgApp.Status.Status, "message", cnvrgApp.Status.Message)
	//	if cnvrgApp.Status.Status == status && cnvrgApp.Status.Message == message {
	//		break
	//	}
	//	if statusCheckAttempts == 0 {
	//		cnvrgAppLog.Info("can't verify status update, status checks attempts exceeded")
	//		break
	//	}
	//	statusCheckAttempts--
	//	cnvrgAppLog.V(1).Info("validating status update", "attempts", statusCheckAttempts)
	//	time.Sleep(1 * time.Second)
	//}
}

func (r *CnvrgAppReconciler) syncCnvrgAppSpec(name types.NamespacedName) (bool, error) {

	cnvrgAppLog.Info("synchronizing cnvrgApp spec")

	// Fetch current cnvrgApp spec
	cnvrgApp, err := r.getCnvrgAppSpec(name)
	if err != nil {
		return false, err
	}
	if cnvrgApp == nil {
		return false, nil // probably cnvrgapp was removed
	}
	cnvrgAppLog = r.Log.WithValues("name", name, "ns", cnvrgApp.Namespace)

	// Get default cnvrgApp spec
	desiredSpec := mlopsv1.DefaultCnvrgAppSpec()

	// Merge current cnvrgApp spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, cnvrgApp.Spec, mergo.WithOverride); err != nil {
		cnvrgAppLog.Error(err, "can't merge")
		return false, err
	}

	equal := reflect.DeepEqual(desiredSpec, cnvrgApp.Spec)
	if !equal {
		cnvrgAppLog.Info("states are not equals, syncing and requeuing")
		cnvrgApp.Spec = desiredSpec
		if err := r.Update(context.Background(), cnvrgApp); err != nil && errors.IsConflict(err) {
			cnvrgAppLog.Info("conflict updating cnvrgApp object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	cnvrgAppLog.Info("states are equals, no need to sync")
	return equal, nil
}

func (r *CnvrgAppReconciler) getCnvrgAppSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgApp, error) {
	ctx := context.Background()
	var cnvrgApp mlopsv1.CnvrgApp
	if err := r.Get(ctx, namespacedName, &cnvrgApp); err != nil {
		if errors.IsNotFound(err) {
			cnvrgAppLog.Info("unable to fetch CnvrgApp, probably cr was deleted")
			return nil, nil
		}
		cnvrgAppLog.Error(err, "unable to fetch CnvrgApp")
		return nil, err
	}
	return &cnvrgApp, nil
}

func (r *CnvrgAppReconciler) cleanup(cnvrgApp *mlopsv1.CnvrgApp) error {

	cnvrgAppLog.Info("running finalizer cleanup")

	// remove cnvrg-db-init
	if err := r.cleanupDbInitCm(cnvrgApp); err != nil {
		return err
	}

	// update infra reconciler cm
	if err := r.triggerInfraReconciler(cnvrgApp, "remove"); err != nil {
		return err
	}

	return nil
}

func (r *CnvrgAppReconciler) cleanupDbInitCm(desiredSpec *mlopsv1.CnvrgApp) error {
	cnvrgAppLog.Info("running cnvrg-db-init cleanup")
	ctx := context.Background()
	dbInitCm := &v1core.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cnvrg-db-init", Namespace: desiredSpec.Namespace}}
	err := r.Delete(ctx, dbInitCm)
	if err != nil && errors.IsNotFound(err) {
		cnvrgAppLog.Info("no need to delete cnvrg-db-init, cm not found")
	} else {
		cnvrgAppLog.Error(err, "error deleting cnvrg-db-init")
		return err
	}
	return nil
}

func (r *CnvrgAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	cnvrgAppLog = r.Log.WithValues("initializing", "crds")

	p := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {

			if reflect.TypeOf(&mlopsv1.CnvrgApp{}) == reflect.TypeOf(e.ObjectOld) {
				oldObject := e.ObjectOld.(*mlopsv1.CnvrgApp)
				newObject := e.ObjectNew.(*mlopsv1.CnvrgApp)
				// deleting cnvrg cr
				if !newObject.ObjectMeta.DeletionTimestamp.IsZero() {
					return true
				}
				shouldReconcileOnSpecChange := reflect.DeepEqual(oldObject.Spec, newObject.Spec) // cnvrgapp spec wasn't changed, assuming status update, won't reconcile
				cnvrgAppLog.V(1).Info("update received", "shouldReconcileOnSpecChange", shouldReconcileOnSpecChange)

				return !shouldReconcileOnSpecChange
			}
			return true
		},
	}

	cnvrgAppController := ctrl.
		NewControllerManagedBy(mgr).
		For(&mlopsv1.CnvrgApp{}).
		WithEventFilter(p)

	for _, v := range desired.Kinds {

		if strings.Contains(v.Group, "istio.io") && viper.GetBool("own-istio-resources") == false {
			continue
		}
		if strings.Contains(v.Group, "openshift.io") && viper.GetBool("own-openshift-resources") == false {
			continue
		}
		if strings.Contains(v.Group, "coreos.com") && viper.GetBool("own-prometheus-resources") == false {
			continue
		}
		u := &unstructured.Unstructured{}
		u.SetGroupVersionKind(v)
		cnvrgAppController.Owns(u)
	}
	cnvrgAppLog.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))
	return cnvrgAppController.
		WithOptions(controller.Options{MaxConcurrentReconciles: viper.GetInt("max-concurrent-reconciles")}).
		Complete(r)
}

func (r *CnvrgAppReconciler) CheckDeploymentReadiness(name types.NamespacedName) (bool, error) {
	ctx := context.Background()
	deployment := &v1apps.Deployment{}

	if err := r.Get(ctx, name, deployment); err != nil && errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if deployment.Status.Replicas == deployment.Status.ReadyReplicas {
		return true, nil
	}

	return false, nil
}

func (r *CnvrgAppReconciler) CheckStatefulSetReadiness(name types.NamespacedName) (bool, error) {

	ctx := context.Background()
	sts := &v1apps.StatefulSet{}

	if err := r.Get(ctx, name, sts); err != nil && errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if sts.Status.Replicas == sts.Status.ReadyReplicas {
		return true, nil
	}

	return false, nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Dimss/crypt/apr1_crypt"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/controlplane"
	"github.com/cnvrg-operator/pkg/dbs"
	"github.com/cnvrg-operator/pkg/desired"
	"github.com/cnvrg-operator/pkg/gpu"
	"github.com/cnvrg-operator/pkg/logging"
	"github.com/cnvrg-operator/pkg/monitoring"
	"github.com/cnvrg-operator/pkg/networking"
	"github.com/cnvrg-operator/pkg/registry"
	"github.com/cnvrg-operator/pkg/shared"
	"github.com/cnvrg-operator/pkg/storage"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/d4l3k/messagediff.v1"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"os"
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

const CnvrginfraFinalizer = "cnvrginfra.mlops.cnvrg.io/finalizer"

type CnvrgInfraReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var cnvrgInfraLog logr.Logger

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrginfras,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrginfras/status,verbs=get;update;patch

func (r *CnvrgInfraReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	cnvrgInfraLog = r.Log.WithValues("name", req.NamespacedName)
	cnvrgInfraLog.Info("starting cnvrginfra reconciliation")

	equal, err := r.syncCnvrgInfraSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !equal {
		requeueAfter, _ := time.ParseDuration("3s")
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	cnvrgInfra, err := r.getCnvrgInfraSpec(req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cnvrgInfra == nil {
		return ctrl.Result{}, nil // probably spec was deleted, no need to reconcile
	}

	// setup finalizer
	if cnvrgInfra.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(cnvrgInfra.ObjectMeta.Finalizers, CnvrginfraFinalizer) {
			cnvrgInfra.ObjectMeta.Finalizers = append(cnvrgInfra.ObjectMeta.Finalizers, CnvrginfraFinalizer)
			if err := r.Update(context.Background(), cnvrgInfra); err != nil {
				cnvrgInfraLog.Error(err, "failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(cnvrgInfra.ObjectMeta.Finalizers, CnvrginfraFinalizer) {
			r.updateStatusMessage(mlopsv1.StatusRemoving, "removing cnvrg spec", cnvrgInfra)
			if err := r.cleanup(cnvrgInfra); err != nil {
				return ctrl.Result{}, err
			}
			cnvrgInfra, err := r.getCnvrgInfraSpec(req.NamespacedName)
			if err != nil {
				return ctrl.Result{}, err
			}
			if cnvrgInfra == nil {
				return ctrl.Result{}, nil
			}
			cnvrgInfra.ObjectMeta.Finalizers = removeString(cnvrgInfra.ObjectMeta.Finalizers, CnvrginfraFinalizer)
			if err := r.Update(context.Background(), cnvrgInfra); err != nil {
				cnvrgInfraLog.Info("error in removing finalizer, checking if cnvrgInfra object still exists")
				//// if update was failed, make sure that cnvrgInfra still exists
				//spec, e := r.getCnvrgInfraSpec(req.NamespacedName)
				//if spec == nil && e == nil {
				//	return ctrl.Result{}, nil // probably spec was deleted, stop reconcile
				//}
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	r.updateStatusMessage(mlopsv1.StatusReconciling, "reconciling", cnvrgInfra)

	//// generate SSO tokens and basic auth secrets
	//if cnvrgInfra.Spec.SSO.Enabled == "true" {
	//	basicAuthUser := "cnvrg"
	//	pubkey, basicAuthPass, err := r.generateSSOKeysAndToken(cnvrgInfra)
	//	if err != nil {
	//		return ctrl.Result{}, err
	//	}
	//	// crate basic auth secret for prometheus
	//	if err := r.createBasicAuthSecret(
	//		cnvrgInfra.Spec.Monitoring.Prometheus.BasicAuthRef,
	//		cnvrgInfra.Spec.InfraNamespace,
	//		basicAuthUser,
	//		string(basicAuthPass),
	//		string(pubkey),
	//		cnvrgInfra,
	//	); err != nil {
	//		return ctrl.Result{}, err
	//	}
	//}

	// apply manifests
	if err := r.applyManifests(cnvrgInfra); err != nil {
		return ctrl.Result{}, err
	}

	// infra reconciler trigger configmap
	if err := r.createInfraReconcilerTriggerCm(cnvrgInfra); err != nil {
		return ctrl.Result{}, err
	}

	r.updateStatusMessage(mlopsv1.StatusHealthy, "successfully reconciled", cnvrgInfra)
	cnvrgInfraLog.Info("successfully reconciled")
	return ctrl.Result{}, nil
}

func (r *CnvrgInfraReconciler) getBasicAuthCreds(secretName string, ns string) (string, string, error) {
	basicAuthCreds := v1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: ns}}
	namespacedName := types.NamespacedName{Namespace: ns, Name: secretName}
	if err := r.Get(context.Background(), namespacedName, &basicAuthCreds); err != nil && errors.IsNotFound(err) {
		cnvrgInfraLog.Info("returning empty token, secret not found", "secret", secretName)
		return "", "", err
	} else if err != nil {
		cnvrgInfraLog.Error(err, "can't get secret", "secret", secretName)
		return "", "", err
	}
	if _, ok := basicAuthCreds.Data["user"]; !ok {
		cnvrgInfraLog.Info("user is missing!")
		return "", "", fmt.Errorf("basic auth pass is missing in secret: %s", secretName)
	}
	if _, ok := basicAuthCreds.Data["pass"]; !ok {
		cnvrgInfraLog.Info("pass is missing!")
		return "", "", fmt.Errorf("basic auth pass is missing in secret: %s", secretName)
	}
	return string(basicAuthCreds.Data["user"]), string(basicAuthCreds.Data["pass"]), nil
}

func (r *CnvrgInfraReconciler) createBasicAuthSecret(name string, ns string, user string, pass string, pubkey string, infra *mlopsv1.CnvrgInfra) error {
	basicAuthCreds := v1.Secret{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns}}
	namespacedName := types.NamespacedName{Namespace: ns, Name: name}

	if err := r.Get(context.Background(), namespacedName, &basicAuthCreds); err != nil && errors.IsNotFound(err) {
		if err := ctrl.SetControllerReference(infra, &basicAuthCreds, r.Scheme); err != nil {
			cnvrgInfraLog.Error(err, "error set controller reference", "name", name)
			return err
		}
		basicAuthCreds.Data = map[string][]byte{"user": []byte(user), "pass": []byte(pass), "pubkey": []byte(pubkey)}
		if err := r.Create(context.Background(), &basicAuthCreds); err != nil {
			cnvrgInfraLog.Error(err, "failed to create secret", "secret", name)
			return err
		}
		return nil
	} else if err != nil {
		cnvrgInfraLog.Error(err, "can't check if secret exists", "secret", name)
		return err
	}

	return nil
}
func (r *CnvrgInfraReconciler) promCredsSecret(infra *mlopsv1.CnvrgInfra) (user string, pass string, err error) {
	user = "cnvrg"
	namespacedName := types.NamespacedName{Name: infra.Spec.Monitoring.Prometheus.CredsRef, Namespace: infra.Spec.InfraNamespace}
	creds := v1core.Secret{ObjectMeta: metav1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := r.Get(context.Background(), namespacedName, &creds); err != nil && errors.IsNotFound(err) {
		if err := ctrl.SetControllerReference(infra, &creds, r.Scheme); err != nil {
			cnvrgAppLog.Error(err, "error set controller reference", "name", namespacedName.Name)
			return "", "", err
		}
		b := make([]byte, 12)
		_, err = rand.Read(b)
		if err != nil {
			cnvrgAppLog.Error(err, "error generating prometheus password")
			return "", "", err
		}
		pass = base64.StdEncoding.EncodeToString(b)
		passHash, err := apr1_crypt.New().Generate([]byte(pass), nil)
		if err != nil {
			cnvrgInfraLog.Error(err, "error generating prometheus hash ")
			return "", "", err
		}
		creds.Data = map[string][]byte{
			"CNVRG_PROMETHEUS_USER": []byte(user),
			"CNVRG_PROMETHEUS_PASS": []byte(pass),
			"htpasswd":              []byte(fmt.Sprintf("%s:%s", user, passHash)),
		}
		if err := r.Create(context.Background(), &creds); err != nil {
			cnvrgInfraLog.Error(err, "error creating prometheus creds", "name", namespacedName.Name)
			return "", "", err
		}

		return user, pass, nil
	} else if err != nil {
		cnvrgInfraLog.Error(err, "can't check if prometheus creds secret exists", "name", namespacedName.Name)
		return "", "", err
	}

	if _, ok := creds.Data["CNVRG_PROMETHEUS_USER"]; !ok {
		err := fmt.Errorf("prometheus creds secret %s missing require field CNVRG_PROMETHEUS_USER", namespacedName.Name)
		cnvrgAppLog.Error(err, "missing required field")
		return "", "", err
	}

	if _, ok := creds.Data["CNVRG_PROMETHEUS_PASS"]; !ok {
		err := fmt.Errorf("prometheus creds secret %s missing require field CNVRG_PROMETHEUS_PASS", namespacedName.Name)
		cnvrgAppLog.Error(err, "missing required field")
		return "", "", err
	}

	return string(creds.Data["CNVRG_PROMETHEUS_USER"]), string(creds.Data["CNVRG_PROMETHEUS_PASS"]), nil

}

func (r *CnvrgInfraReconciler) generateSSOKeysAndToken(infra *mlopsv1.CnvrgInfra) ([]byte, []byte, error) {

	keysForTokens := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: infra.Spec.SSO.KeysSecretRef, Namespace: infra.Spec.InfraNamespace},
	}
	namespacedName := types.NamespacedName{
		Namespace: keysForTokens.ObjectMeta.Namespace,
		Name:      keysForTokens.ObjectMeta.Name,
	}

	if err := r.Get(context.Background(), namespacedName, &keysForTokens); err != nil && errors.IsNotFound(err) {
		cnvrgInfraLog.Info("secret does not exists, will create one", "secret", infra.Spec.SSO.KeysSecretRef)
		privateKey, err := shared.GeneratePrivateKey()
		if err != nil {
			cnvrgInfraLog.Error(err, "error generating private key for SSO tokens")
			return nil, nil, err
		}
		privateKeyBytes, publicKeyBytes, err := shared.EncodeKeysToPEM(privateKey)
		if err != nil {
			cnvrgInfraLog.Error(err, "error generating public key or encoding keys for SSO tokens")
			return nil, nil, err
		}
		ssoToken, err := shared.CreateSSOToken(infra.Spec.SSO.ClientID, privateKeyBytes)
		if err != nil {
			cnvrgInfraLog.Error(err, "error generating SSO Token")
			return nil, nil, err
		}
		if err := ctrl.SetControllerReference(infra, &keysForTokens, r.Scheme); err != nil {
			cnvrgInfraLog.Error(err, "error set controller reference")
			return nil, nil, err
		}
		keysForTokens.Data = map[string][]byte{
			"pubkey":   publicKeyBytes,
			"privkey":  privateKeyBytes,
			"ssoToken": []byte(ssoToken),
		}
		if err := r.Create(context.Background(), &keysForTokens); err != nil {
			cnvrgInfraLog.Error(err, "failed to create secret for sso tokens ")
			return nil, nil, err
		}
		return publicKeyBytes, []byte(ssoToken), nil
	} else if err != nil {
		cnvrgInfraLog.Error(err, "can't check if secret exists", "secret", infra.Spec.SSO.KeysSecretRef)
		return nil, nil, err
	}
	cnvrgInfraLog.Info("secret exists, skipping...", "secret", infra.Spec.SSO.KeysSecretRef)

	if _, ok := keysForTokens.Data["ssoToken"]; !ok {
		cnvrgInfraLog.Info("sso token is missing!")
		return nil, nil, fmt.Errorf("sso token is missing in secret: %s", infra.Spec.SSO.KeysSecretRef)
	}
	if _, ok := keysForTokens.Data["pubkey"]; !ok {
		cnvrgInfraLog.Info("pubkey token is missing!")
		return nil, nil, fmt.Errorf("pubkey is missing in secret: %s", infra.Spec.SSO.KeysSecretRef)
	}
	return keysForTokens.Data["pubkey"], keysForTokens.Data["ssoToken"], nil

}

func (r *CnvrgInfraReconciler) getCnvrgAppInstances(infra *mlopsv1.CnvrgInfra) ([]mlopsv1.AppInstance, error) {

	cmName := types.NamespacedName{Namespace: infra.Spec.InfraNamespace, Name: infra.Spec.InfraReconcilerCm}
	if cmName.Name == "" {
		cmName.Name = mlopsv1.DefaultCnvrgInfraSpec().InfraReconcilerCm
	}
	cnvrgAppCm := &v1.ConfigMap{}
	if err := r.Get(context.Background(), cmName, cnvrgAppCm); err != nil && errors.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var apps []mlopsv1.AppInstance
	for _, appJson := range cnvrgAppCm.Data {
		var app mlopsv1.AppInstance
		if err := json.Unmarshal([]byte(appJson), &app); err != nil {
			cnvrgInfraLog.Error(err, "error decoding AppInstance")
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (r *CnvrgInfraReconciler) applyManifests(cnvrgInfra *mlopsv1.CnvrgInfra) error {

	var reconcileResult error

	// registry
	cnvrgInfraLog.Info("applying registry")
	registryData := desired.TemplateData{
		Namespace: cnvrgInfra.Spec.InfraNamespace,
		Data:      map[string]interface{}{"Registry": cnvrgInfra.Spec.Registry},
	}
	if err := desired.Apply(registry.State(registryData), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// logging
	cnvrgInfraLog.Info("applying logging")
	cnvrgApps, err := r.getCnvrgAppInstances(cnvrgInfra)
	if err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}
	fluentbitData := desired.TemplateData{
		Namespace: cnvrgInfra.Spec.InfraNamespace,
		Data:      map[string]interface{}{"AppInstance": cnvrgApps},
	}
	if err := desired.Apply(logging.FluentbitConfigurationState(fluentbitData), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}
	if err := desired.Apply(logging.InfraLoggingState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// storage
	cnvrgInfraLog.Info("applying storage")
	if err := desired.Apply(storage.State(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// grafana dashboards
	cnvrgInfraLog.Info("applying grafana dashboards")
	if err := r.createGrafanaDashboards(cnvrgInfra); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// grafana datasource
	cnvrgInfraLog.Info("applying grafana datasource")
	basicAuthUser, basicAuthPass, err := r.promCredsSecret(cnvrgInfra)
	if err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}
	grafanaDatasourceData := desired.TemplateData{
		Namespace: cnvrgInfra.Spec.InfraNamespace,
		Data: map[string]interface{}{
			"Svc":  cnvrgInfra.Spec.Monitoring.Prometheus.SvcName,
			"Port": cnvrgInfra.Spec.Monitoring.Prometheus.Port,
			"User": basicAuthUser,
			"Pass": basicAuthPass,
		},
	}
	if err := desired.Apply(monitoring.GrafanaDSState(grafanaDatasourceData), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	_, _, err = r.promCredsSecret(cnvrgInfra)
	if err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}
	// monitoring
	cnvrgInfraLog.Info("applying monitoring")
	if err := desired.Apply(monitoring.InfraMonitoringState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// istio
	cnvrgInfraLog.Info("applying istio")
	if err := desired.Apply(networking.IstioInstanceState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// redis
	cnvrgInfraLog.Info("applying redis")
	if err := desired.Apply(dbs.InfraDbsState(cnvrgInfra), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// mpi infra
	cnvrgInfraLog.Info("applying mpi infra")
	if err := desired.Apply(controlplane.MpiInfraState(), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
		r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
		reconcileResult = err
	}

	// nvidia device plugin
	if cnvrgInfra.Spec.Gpu.NvidiaDp.Enabled == "true" {
		cnvrgInfraLog.Info("nvidia device plugin")
		nvidiaDpData := desired.TemplateData{
			Namespace: cnvrgInfra.Spec.InfraNamespace,
			Data: map[string]interface{}{
				"NvidiaDp": cnvrgInfra.Spec.Gpu.NvidiaDp,
				"Registry": cnvrgInfra.Spec.Registry,
			},
		}
		if err := desired.Apply(gpu.NvidiaDpState(nvidiaDpData), cnvrgInfra, r.Client, r.Scheme, cnvrgInfraLog); err != nil {
			r.updateStatusMessage(mlopsv1.StatusError, err.Error(), cnvrgInfra)
			reconcileResult = err
		}
	}

	return reconcileResult
}

func (r *CnvrgInfraReconciler) createGrafanaDashboards(cnvrgInfra *mlopsv1.CnvrgInfra) error {

	if cnvrgInfra.Spec.Monitoring.Enabled != "true" {
		cnvrgInfraLog.Info("monitoring disabled, skipping grafana deployment")
		return nil
	}

	basePath := "/pkg/monitoring/tmpl/grafana/dashboards-data/"
	for _, dashboard := range desired.GrafanaInfraDashboards {
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
				Namespace: cnvrgInfra.Spec.InfraNamespace,
			},
			Data: map[string]string{filepath.Base(f.Name()): string(b)},
		}
		if err := ctrl.SetControllerReference(cnvrgInfra, cm, r.Scheme); err != nil {
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

func (r *CnvrgInfraReconciler) syncCnvrgInfraSpec(name types.NamespacedName) (bool, error) {

	cnvrgInfraLog.Info("synchronizing cnvrgInfra spec")

	// Fetch current cnvrgInfra spec
	cnvrgInfra, err := r.getCnvrgInfraSpec(name)
	if err != nil {
		return false, err
	}
	if cnvrgInfra == nil {
		return false, nil // probably cnvrgapp was removed
	}
	cnvrgInfraLog = r.Log.WithValues("name", name, "ns", cnvrgInfra.Spec.InfraNamespace)

	// Get default cnvrgInfra spec
	desiredSpec := mlopsv1.DefaultCnvrgInfraSpec()

	// Merge current cnvrgInfra spec into default spec ( make it indeed desiredSpec )
	if err := mergo.Merge(&desiredSpec, cnvrgInfra.Spec, mergo.WithOverride); err != nil {
		cnvrgInfraLog.Error(err, "can't merge")
		return false, err
	}

	if viper.GetBool("verbose") {
		cnvrgInfraLog.V(1).Info("printing the diff between desiredSpec and actual")
		diff, _ := messagediff.PrettyDiff(desiredSpec, cnvrgInfra.Spec)
		cnvrgInfraLog.V(1).Info(diff)
	}

	// Compare desiredSpec and current cnvrgInfra spec,
	// if they are not equal, update the cnvrgInfra spec with desiredSpec,
	// and return true for triggering new reconciliation
	equal := reflect.DeepEqual(desiredSpec, cnvrgInfra.Spec)
	if !equal {
		cnvrgInfraLog.Info("states are not equals, syncing and requeuing")
		cnvrgInfra.Spec = desiredSpec
		if err := r.Update(context.Background(), cnvrgInfra); err != nil && errors.IsConflict(err) {
			cnvrgAppLog.Error(err, "conflict updating cnvrgInfra object, requeue for reconciliations...")
			return true, nil
		} else if err != nil {
			return false, err
		}
		return equal, nil
	}

	cnvrgInfraLog.Info("states are equals, no need to sync")
	return equal, nil
}

func (r *CnvrgInfraReconciler) getCnvrgInfraSpec(namespacedName types.NamespacedName) (*mlopsv1.CnvrgInfra, error) {
	ctx := context.Background()
	var cnvrgInfra mlopsv1.CnvrgInfra
	if err := r.Get(ctx, namespacedName, &cnvrgInfra); err != nil {
		if errors.IsNotFound(err) {
			cnvrgInfraLog.Info("unable to fetch CnvrgApp, probably cr was deleted")
			return nil, nil
		}
		cnvrgInfraLog.Error(err, "unable to fetch CnvrgApp")
		return nil, err
	}
	return &cnvrgInfra, nil
}

func (r *CnvrgInfraReconciler) cleanup(cnvrgInfra *mlopsv1.CnvrgInfra) error {

	cnvrgInfraLog.Info("running finalizer cleanup")

	// remove istio
	if err := r.cleanupIstio(cnvrgInfra); err != nil {
		return err
	}

	// cleanup pvc
	if err := r.cleanupPVCs(cnvrgInfra); err != nil {
		return err
	}

	// cleanup tokens keys cm
	if err := r.cleanupTokensKeysCM(cnvrgInfra); err != nil {
		return err
	}

	return nil
}

func (r *CnvrgInfraReconciler) cleanupTokensKeysCM(infra *mlopsv1.CnvrgInfra) error {
	cnvrgAppLog.Info("running cleanup", "secret", infra.Spec.SSO.KeysSecretRef)
	ctx := context.Background()
	tokensKeySecret := &v1core.Secret{ObjectMeta: metav1.ObjectMeta{Name: infra.Spec.SSO.KeysSecretRef, Namespace: infra.Spec.InfraNamespace}}
	err := r.Delete(ctx, tokensKeySecret)
	if err != nil && errors.IsNotFound(err) {
		cnvrgAppLog.Info("no need to delete secret ", "secret", infra.Spec.SSO.KeysSecretRef)
	} else if err != nil {
		cnvrgAppLog.Error(err, "error deleting ", "secret", infra.Spec.SSO.KeysSecretRef)
		return err
	}
	return nil
}

func (r *CnvrgInfraReconciler) cleanupPVCs(infra *mlopsv1.CnvrgInfra) error {
	cnvrgAppLog.Info("running pvc cleanup")
	ctx := context.Background()
	pvcList := v1core.PersistentVolumeClaimList{}
	if err := r.List(ctx, &pvcList); err != nil {
		cnvrgAppLog.Error(err, "failed cleanup pvcs")
		return err
	}
	for _, pvc := range pvcList.Items {
		if pvc.Namespace == infra.Spec.InfraNamespace {
			if _, ok := pvc.ObjectMeta.Labels["app"]; ok {
				if pvc.ObjectMeta.Labels["app"] == "prometheus" {
					if err := r.Delete(ctx, &pvc); err != nil && errors.IsNotFound(err) {
						cnvrgInfraLog.Info("prometheus pvc already deleted")
					} else if err != nil {
						cnvrgInfraLog.Error(err, "error deleting prometheus pvc")
						return err
					}
				}
			}
		}
	}
	return nil
}

func (r *CnvrgInfraReconciler) cleanupIstio(cnvrgInfra *mlopsv1.CnvrgInfra) error {
	cnvrgInfraLog.Info("running istio cleanup")
	ctx := context.Background()
	istioManifests := networking.IstioInstanceState(cnvrgInfra)
	for _, m := range istioManifests {
		// Make sure IstioOperator was deployed
		if m.GVR == desired.Kinds[desired.IstioGVR] {
			if m.TemplateData == nil {
				m.TemplateData = cnvrgInfra
			}
			if err := m.GenerateDeployable(); err != nil {
				cnvrgInfraLog.Error(err, "can't make manifest deployable")
				return err
			}
			if err := r.Delete(ctx, m.Obj); err != nil {
				if errors.IsNotFound(err) {
					cnvrgInfraLog.Info("istio instance not found - probably removed previously")
					return nil
				}
				return err
			}
			istioExists := true
			cnvrgInfraLog.Info("wait for istio instance removal")
			for istioExists {
				err := r.Get(ctx, types.NamespacedName{Name: m.Obj.GetName(), Namespace: m.Obj.GetNamespace()}, m.Obj)
				if err != nil && errors.IsNotFound(err) {
					cnvrgInfraLog.Info("istio instance was successfully removed")
					istioExists = false
				}
				if istioExists {
					cnvrgInfraLog.Info("istio instance still present, will sleep of 1 sec, and check again...")
				}
			}
		}
	}
	return nil
}

func (r *CnvrgInfraReconciler) updateStatusMessage(status mlopsv1.OperatorStatus, message string, cnvrgInfra *mlopsv1.CnvrgInfra) {
	if cnvrgInfra.Status.Status == mlopsv1.StatusRemoving {
		cnvrgInfraLog.Info("skipping status update, current cnvrg spec under removing status...")
		return
	}
	ctx := context.Background()
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		name := types.NamespacedName{Namespace: "", Name: cnvrgInfra.Name}
		infra, err := r.getCnvrgInfraSpec(name)
		if err != nil {
			return err
		}
		infra.Status.Status = status
		infra.Status.Message = message
		err = r.Status().Update(ctx, infra)
		return err
	})
	if err != nil {
		cnvrgInfraLog.Error(err, "can't update status")
	}
}

func (r *CnvrgInfraReconciler) createInfraReconcilerTriggerCm(cnvrgInfra *mlopsv1.CnvrgInfra) error {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: cnvrgInfra.Spec.InfraReconcilerCm, Namespace: cnvrgInfra.Spec.InfraNamespace},
	}
	if err := ctrl.SetControllerReference(cnvrgInfra, cm, r.Scheme); err != nil {
		cnvrgInfraLog.Error(err, "failed to set ControllerReference", "cm", cnvrgInfra.Spec.InfraReconcilerCm)
		return err
	}
	if err := r.Create(context.Background(), cm); err != nil && errors.IsAlreadyExists(err) {
		cnvrgInfraLog.Info("already exists", "cm", cnvrgInfra.Spec.InfraReconcilerCm)
	} else if err != nil {
		cnvrgInfraLog.Error(err, "error creating", "cm", cnvrgInfra.Spec.InfraReconcilerCm)
		return err
	}

	return nil
}

func (r *CnvrgInfraReconciler) SetupWithManager(mgr ctrl.Manager) error {
	cnvrgInfraLog = r.Log.WithValues("initializing", "crds")

	if viper.GetBool("deploy-depended-crds") == false {
		zap.S().Info("deploy-depended-crds is false, I hope CRDs was deployed ahead and match expected versions, if not I will fail...")
	} else {

		if viper.GetBool("own-istio-resources") {
			err := desired.Apply(networking.IstioCrds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r, r.Scheme, r.Log)
			if err != nil {
				cnvrgInfraLog.Error(err, "can't apply istio CRDs")
				os.Exit(1)
			}
		}

		if viper.GetBool("own-prometheus-resources") {
			err := desired.Apply(monitoring.Crds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r, r.Scheme, r.Log)
			if err != nil {
				cnvrgInfraLog.Error(err, "can't apply prometheus CRDs")
				os.Exit(1)
			}
		}

		err := desired.Apply(controlplane.Crds(), &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}, r, r.Scheme, r.Log)
		if err != nil {
			cnvrgInfraLog.Error(err, "can't apply MPI CRDs")
			os.Exit(1)
		}
	}

	p := predicate.Funcs{

		UpdateFunc: func(e event.UpdateEvent) bool {

			// run reconcile only changing cnvrginfra/object marked for deletion
			if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(e.ObjectOld) {
				oldObject := e.ObjectOld.(*mlopsv1.CnvrgInfra)
				newObject := e.ObjectNew.(*mlopsv1.CnvrgInfra)
				// deleting cnvrg cr
				if !newObject.ObjectMeta.DeletionTimestamp.IsZero() {
					return true
				}
				shouldReconcileOnSpecChange := reflect.DeepEqual(oldObject.Spec, newObject.Spec) // cnvrginfra spec wasn't changed, assuming status update, won't reconcile
				cnvrgInfraLog.V(1).Info("cnvrginfra update received", "shouldReconcileOnSpecChange", shouldReconcileOnSpecChange)

				return !shouldReconcileOnSpecChange

			}
			return true
		},
	}

	cnvrgInfraController := ctrl.
		NewControllerManagedBy(mgr).
		For(&mlopsv1.CnvrgInfra{}).
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
		cnvrgInfraController.Owns(u)
	}
	cnvrgInfraLog.Info(fmt.Sprintf("max concurrent reconciles: %d", viper.GetInt("max-concurrent-reconciles")))
	return cnvrgInfraController.
		WithOptions(controller.Options{MaxConcurrentReconciles: viper.GetInt("max-concurrent-reconciles")}).
		Complete(r)
}

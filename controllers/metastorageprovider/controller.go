package metastorageprovider

import (
	"context"
	"fmt"
	v1alpha1 "github.com/AccessibleAI/cnvrg-operator/api/v1alpha1"
	installer "github.com/AccessibleAI/cnvrg-operator/pkg/helm"
	v1apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	finalizer                 = "metastorageprovisioner.finalizers.metacloud.io"
	nfsProvisionerReleaseName = "nfs-subdir-external-provisioner"
	nfsProvisionerVersion     = "4.0.2"
)

// Reconciler reconciles a MetaStorageProvisioner object
type Reconciler struct {
	Client      client.Client
	recorder    record.EventRecorder
	Scheme      *runtime.Scheme
	Namespace   string
	EventLogger EventLogger
}

// readinessChecker is a function that checks the readiness of a provisioner workload
type readinessChecker func() (bool, error)

// +kubebuilder:rbac:groups=mlops.cnvrg.io,namespace=cnvrg,resources=metastorageprovisioners,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,namespace=cnvrg,resources=metastorageprovisioners/status,verbs=get;update;patch
// Reconcile reconciles the MetaStorageProvisioner resource
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	r.EventLogger = EventLogger{logger: logger, recorder: r.recorder, client: r.Client}

	r.EventLogger.WithMessage("reconciling storage provisioner").Log(ctx)

	provisionerObject := &v1alpha1.MetaStorageProvisioner{}

	// get reconciled object
	if err := r.Client.Get(ctx, req.NamespacedName, provisionerObject); err != nil {
		if errors.IsNotFound(err) {
			r.EventLogger.WithMessage("reconciled meta storage object not found - ignoring").Log(ctx)
			return ctrl.Result{}, nil
		}
		r.EventLogger.WithError(err).WithMessage("failed to get meta storage object").Log(ctx)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	r.EventLogger.Object = provisionerObject

	if provisionerObject.GetDeletionTimestamp().IsZero() {
		if err := r.addFinalizer(ctx, provisionerObject); err != nil {
			r.EventLogger.WithError(err).WithMessage("failed to add finalizer").Log(ctx)
			return ctrl.Result{}, err
		}
	} else {
		if err := r.ReconcileDelete(ctx, provisionerObject); err != nil {
			r.EventLogger.WithError(err).WithMessage("failed to reconcile delete").WithStatus(v1alpha1.Failed).Log(ctx)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// validate the provisioner object
	if err := provisionerObject.Spec.Validate(); err != nil {
		r.EventLogger.WithError(err).WithMessage("failed to validate provisioner object").WithStatus(v1alpha1.Failed).Log(ctx)
		return ctrl.Result{}, err
	}

	var chartConfig installer.ChartConfig
	var checker readinessChecker
	values := make(map[string]interface{})

	if provisionerObject.Spec.NFSProvisioner != nil {
		// NFS provisioner chart configuration
		chartConfig = r.nfsConfig(provisionerObject, values, chartConfig)

		// readiness checker for NFS provisioner
		checker = func() (bool, error) {
			deployment := types.NamespacedName{Name: provisionerObject.Name, Namespace: r.Namespace}
			return CheckDeploymentReadiness(r.Client, deployment)
		}
	}

	// install the provisioner
	helm, err := installer.NewHelm(chartConfig, logger)
	if err != nil {
		r.EventLogger.WithError(err).WithMessage("failed to create helm installer").WithStatus(v1alpha1.Failed).Log(ctx)
		return ctrl.Result{}, err
	}

	if err := helm.Install(); err != nil {
		r.EventLogger.WithError(err).WithMessage("failed to install provisioner").WithStatus(v1alpha1.Failed).Log(ctx)
		return ctrl.Result{}, err
	}

	// check readiness of the provisioner
	var ready bool
	ready, err = checker()
	if err != nil {
		r.EventLogger.WithError(err).WithMessage("failed to check readiness").WithStatus(v1alpha1.Failed).Log(ctx)
		return ctrl.Result{}, err
	}

	// if not ready, requeue
	if !ready {
		r.EventLogger.WithMessage("waiting for deployment readiness").WithStatus(v1alpha1.Pending).Log(ctx)
		return ctrl.Result{RequeueAfter: 30}, nil
	}

	// provisioner is ready
	r.EventLogger.WithMessage("successfully installed provisioner").WithStatus(v1alpha1.Running).Log(ctx)
	return ctrl.Result{}, nil
}

func (r *Reconciler) nfsConfig(provisionerObject *v1alpha1.MetaStorageProvisioner, values map[string]interface{}, chartConfig installer.ChartConfig) installer.ChartConfig {
	var storageClass map[string]interface{}

	if provisionerObject.Spec.NFSProvisioner.StorageClassName != "" {
		storageClass = map[string]interface{}{
			"name": provisionerObject.Spec.NFSProvisioner.StorageClassName,
		}
		values["storageClass"] = storageClass
	}

	values["nfs"] = map[string]interface{}{
		"server": provisionerObject.Spec.NFSProvisioner.NFSServer,
		"path":   provisionerObject.Spec.NFSProvisioner.NFSPath,
	}

	values["fullnameOverride"] = provisionerObject.Name

	chartConfig = r.NFSChartConfig(provisionerObject.Name, values)
	return chartConfig
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.recorder = mgr.GetEventRecorderFor("cnvrg-meta-storage-provisioner")
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		For(&v1alpha1.MetaStorageProvisioner{}).
		Complete(r)
}

func (r *Reconciler) NFSChartConfig(name string, values map[string]interface{}) installer.ChartConfig {
	return installer.ChartConfig{
		Namespace:   r.Namespace,
		ReleaseName: name,
		ChartName:   nfsProvisionerReleaseName,
		Url:         fmt.Sprintf("https://kubernetes-sigs.github.io/%s", nfsProvisionerReleaseName),
		Values:      values,
		Version:     nfsProvisionerVersion,
	}
}

func (r *Reconciler) addFinalizer(ctx context.Context, provisioner *v1alpha1.MetaStorageProvisioner) error {
	// if finalizer not present on Provisioner object, add it
	if !controllerutil.ContainsFinalizer(provisioner, finalizer) {
		controllerutil.AddFinalizer(provisioner, finalizer)
		return r.Client.Update(ctx, provisioner)
	}
	return nil
}

// ReconcileDelete reconciles the deletion of the MetaStorageProvisioner resource
func (r *Reconciler) ReconcileDelete(ctx context.Context, provisioner *v1alpha1.MetaStorageProvisioner) error {
	if !controllerutil.ContainsFinalizer(provisioner, finalizer) {
		return nil
	}

	logger := log.FromContext(ctx)
	r.EventLogger.WithMessage("deleting meta storage provisioner").WithStatus(v1alpha1.Deleting).Log(ctx)

	chartConfig := r.NFSChartConfig(provisioner.Name, nil)
	helm, err := installer.NewHelm(chartConfig, logger)
	if err != nil {
		return fmt.Errorf("error while creating helm installer: %w", err)
	}

	err = helm.Delete()
	if err != nil {
		return fmt.Errorf("error while deleting provisioner: %w", err)
	}

	controllerutil.RemoveFinalizer(provisioner, finalizer)
	if err := r.Client.Update(ctx, provisioner); err != nil {
		return fmt.Errorf("error while removing finalizer from meta storage provisioner object: %w", err)
	}

	return nil
}

func CheckDeploymentReadiness(client client.Client, name types.NamespacedName) (bool, error) {
	ctx := context.Background()
	deployment := &v1apps.Deployment{}

	if err := client.Get(ctx, name, deployment); err != nil && errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if deployment.Status.Replicas == deployment.Status.ReadyReplicas {
		return true, nil
	}

	return false, nil
}

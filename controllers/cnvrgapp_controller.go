/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/pg"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

// CnvrgAppReconciler reconciles a CnvrgApp object
type CnvrgAppReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=cnvrgapps/status,verbs=get;update;patch

func (r *CnvrgAppReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	r.Log.Info("starting reconciliation")
	desiredSpec, err := r.desiredSpec(req)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err := r.pg(desiredSpec); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *CnvrgAppReconciler) desiredSpec(req ctrl.Request) (*mlopsv1.CnvrgApp, error) {
	ctx := context.Background()
	var cnvrgApp mlopsv1.CnvrgApp
	if err := r.Get(ctx, req.NamespacedName, &cnvrgApp); err != nil {
		r.Log.Info("unable to fetch CnvrgApp, probably cr was deleted")
		return nil, nil
	}
	desiredSpec := mlopsv1.CnvrgApp{Spec: mlopsv1.DefaultSpec}
	if err := mergo.Merge(&desiredSpec, cnvrgApp, mergo.WithOverride); err != nil {
		r.Log.Error(err, "can't merge")
		return nil, err
	}
	return &desiredSpec, nil
}

func (r *CnvrgAppReconciler) pg(desiredSpec *mlopsv1.CnvrgApp) error {
	ctx := context.Background()
	for _, s := range pg.State(desiredSpec) {
		if err := s.GenerateDeployable(desiredSpec); err != nil {
			r.Log.Error(err, "error generating deployable", "name", s.Name)
			return err
		}
		if err := ctrl.SetControllerReference(desiredSpec, s.Obj, r.Scheme); err != nil {
			r.Log.Error(err, "error setting controller reference", "name", s.Name)
			return err
		}
		if err := r.Create(ctx, s.Obj); err != nil {
			r.Log.Error(err, "error creating object", "name", s.Name)
			return err
		}
	}
	return nil
}

func (r *CnvrgAppReconciler) SetupWithManager(mgr ctrl.Manager) error {

	deployments := &unstructured.Unstructured{}
	deployments.SetGroupVersionKind(schema.GroupVersionKind{Kind: "Deployment", Group: "", Version: "apps/v1"})

	services := &unstructured.Unstructured{}
	services.SetGroupVersionKind(schema.GroupVersionKind{Kind: "Service", Group: "", Version: "v1"})

	pvcs := &unstructured.Unstructured{}
	pvcs.SetGroupVersionKind(schema.GroupVersionKind{Kind: "PersistentVolumeClaim", Group: "", Version: "v1"})

	secrets := &unstructured.Unstructured{}
	secrets.SetGroupVersionKind(schema.GroupVersionKind{Kind: "Secret", Group: "", Version: "v1"})

	return ctrl.NewControllerManagedBy(mgr).
		For(&mlopsv1.CnvrgApp{}).
		Owns(&corev1.ConfigMap{}).
		Owns(deployments).
		Owns(services).
		Owns(pvcs).
		Owns(secrets).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}

// Helper functions to check and remove string from a slice of strings.
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

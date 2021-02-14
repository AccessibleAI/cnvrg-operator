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
	"fmt"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/cnvrg-operator/pkg/pg"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/markbates/pkger"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"text/template"
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
	ctx := context.Background()
	r.Log.Info("asd")
	_ = r.Log.WithValues("cnvrgapp", req.NamespacedName)
	r.Log.Info("RC Name: " + req.NamespacedName.Name)
	r.Log.Info("This is my first k8s controller!!!!!")
	var cnvrgApp mlopsv1.CnvrgApp
	if err := r.Get(ctx, req.NamespacedName, &cnvrgApp); err != nil {
		r.Log.Info("unable to fetch CnvrgApp, probably cr was deleted")
		return ctrl.Result{}, nil
	}

	//cnvrgFinalizer := "cnvrgapp.finalizers.cnvrg.io"

	//if cnvrgApp.ObjectMeta.DeletionTimestamp.IsZero() {
	//	if !containsString(cnvrgApp.ObjectMeta.Finalizers, cnvrgFinalizer) {
	//		cnvrgApp.ObjectMeta.Finalizers = append(cnvrgApp.ObjectMeta.Finalizers, cnvrgFinalizer)
	//		if err := r.Update(ctx, &cnvrgApp); err != nil {
	//			return ctrl.Result{}, err
	//		}
	//	}
	//} else {
	//	if containsString(cnvrgApp.ObjectMeta.Finalizers, cnvrgFinalizer) {
	//		r.Log.Info("I'm a finalizer, purging some stuff here ....")
	//
	//		cnvrgApp.ObjectMeta.Finalizers = removeString(cnvrgApp.ObjectMeta.Finalizers, cnvrgFinalizer)
	//
	//		if err := r.Update(ctx, &cnvrgApp); err != nil {
	//			return ctrl.Result{}, err
	//		}
	//	}
	//}
	dep := getDep()
	cm := getCm()
	ctrl.SetControllerReference(&cnvrgApp, cm, r.Scheme)
	err := ctrl.SetControllerReference(&cnvrgApp, dep, r.Scheme)
	if err != nil {
		r.Log.Info("This is error")
	}

	r.Create(ctx, dep)

	err = r.Get(ctx, types.NamespacedName{Name: "test-cm", Namespace: "default"}, cm)
	if err != nil && errors.IsNotFound(err) {
		if err := r.Create(ctx, cm); err != nil {
			r.Log.Error(err, "error creating cm")
			return ctrl.Result{}, err
		}
	}
	f, err := pkger.Open("/pkg/pg/tmpl/svc.tpl")
	if err != nil {

	}
	b, err := ioutil.ReadAll(f)
	fmt.Print(string(b))
	tmpl, err := template.New("pg-pvc").Parse(string(b))
	//tmpl, err := template.ParseFiles("/pkg/db/pg/pvc.tpl")
	defaultCnvrgApp := mlopsv1.CnvrgApp{Spec: mlopsv1.DefaultCnvrgAppSpec()}

	if err := mergo.Merge(&defaultCnvrgApp, cnvrgApp, mergo.WithOverride); err != nil {
		// ...
	}
	pg.Deploy()

	err = tmpl.Execute(os.Stdout, defaultCnvrgApp)
	if err != nil {
		r.Log.Error(err, "error parsing template")
	}

	return ctrl.Result{}, nil
}

func (r *CnvrgAppReconciler) manageConfigMap(cnvrgApp *mlopsv1.CnvrgApp) (*reconcile.Result, error) {
	r.Log.Info("This is CM stuff... ")
	return nil, nil

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

func getCm() *corev1.ConfigMap {
	configMapData := make(map[string]string, 0)
	uiProperties := `color.good=purple`
	configMapData["ui.properties"] = uiProperties
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "default",
		},
		Data: configMapData,
	}
	return cm
}

func getDep() *unstructured.Unstructured {
	const deploymentYAML = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: default
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
`
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Kind:    "Deployment",
		Group:   "",
		Version: "apps/v1",
	})

	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, _ := dec.Decode([]byte(deploymentYAML), nil, obj)

	fmt.Println(obj.GetName(), gvk.String())

	return obj

}

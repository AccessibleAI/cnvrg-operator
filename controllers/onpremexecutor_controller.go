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
	"bytes"
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/markbates/pkger"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"text/template"

	mlopsv1 "github.com/cnvrg-operator/api/v1"
)

// OnPremExecutorReconciler reconciles a OnPremExecutor object
type OnPremExecutorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=onpremexecutors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mlops.cnvrg.io,resources=onpremexecutors/status,verbs=get;update;patch

func (r *OnPremExecutorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("onpremexecutor", req.NamespacedName)

	r.Log.Info("This is Cnvrg OnPrem Executor ")
	ctx := context.Background()
	var executor mlopsv1.OnPremExecutor
	if err := r.Get(ctx, req.NamespacedName, &executor); err != nil {
		r.Log.Error(err, "unable to fetch OnPremExecutor")
	}
	envVars := r.jobEnvVars(executor)
	jobScript := generateExecutorScript(envVars)
	r.Log.Info(jobScript)
	if err := dumpJobScript(jobScript); err != nil {
		r.Log.Error(err, "error dumping on-prem-executor.sh script")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *OnPremExecutorReconciler) jobEnvVars(executor mlopsv1.OnPremExecutor) []string {
	r.Log.Info("dumping job env vars into host")
	var envVars []string
	envVar := executor.Spec.Template.Spec.Containers[0].Env
	for i := 0; i < len(envVar); i++ {
		envVars = append(envVars, fmt.Sprintf(`export %s=%s`, envVar[i].Name, envVar[i].Value))
	}
	return envVars
}

func generateExecutorScript(jobEnvVars []string) string {
	templateData := map[string]interface{}{
		"Data": map[string]interface{}{
			"JobEnvVars": jobEnvVars,
		},
	}
	buffer, err := renderTemplate("/pkg/executor/tmpl/exec_cnvrg_job.tpl", templateData)
	if err != nil {
		logrus.Errorf("error generating cluster setup script err: %v", err)
		panic(err)
	}
	executorScript := buffer.String()
	return executorScript
}

func dumpJobScript(jobScript string) error {
	err := ioutil.WriteFile("/tmp/onprem-executor.sh", []byte(jobScript), 0755)
	if err != nil {
		return err
	}

	return nil
}

func renderTemplate(templateFile string, templateData map[string]interface{}) (*bytes.Buffer, error) {
	var tpl bytes.Buffer
	f, err := pkger.Open(templateFile)
	if err != nil {
		logrus.Errorf("error reading template %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		logrus.Errorf("%v, error reading file: %v", err, templateFile)
		return nil, err
	}
	clusterTmpl, err := template.New(strings.ReplaceAll(templateFile, "/", "-")).Parse(string(b))
	if err != nil {
		logrus.Errorf("%v, template: %v", err, templateFile)
		return nil, err
	}
	if err = clusterTmpl.Execute(&tpl, templateData); err != nil {
		logrus.Errorf("err: %v rendering template error", err)
		return nil, err
	}
	return &tpl, nil
}

func (r *OnPremExecutorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mlopsv1.OnPremExecutor{}).
		Complete(r)
}

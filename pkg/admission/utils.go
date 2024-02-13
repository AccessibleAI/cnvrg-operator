package admission

import (
	"github.com/AccessibleAI/cnvrg-shim/apis/metacloud/v1alpha1"
	"go.uber.org/zap"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func endWithError(err error, w http.ResponseWriter) {
	zap.S().Error(err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func endWithOk(data []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		zap.S().Error(err)
	}
}

func KubeClient() client.Client {
	ctrl.SetLogger(ctrlzap.New(ctrlzap.UseFlagOptions(&ctrlzap.Options{Development: false})))
	rc, err := config.GetConfig()
	if err != nil {
		zap.S().Fatal(err)
	}

	scheme := runtime.NewScheme()
	utilruntime.Must(v1core.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	cc, err := client.New(rc, client.Options{Scheme: scheme})
	if err != nil {
		zap.S().Fatal(err)
	}
	return cc
}

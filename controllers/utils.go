package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/networking"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
	"strconv"
	"strings"
)

type cnvrgSpecBoolTransformer struct{}

func (t cnvrgSpecBoolTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(true) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				// always set boolean value
				// e.g always do the WithOverwriteWithEmptyValue
				// but only for booleans
				dst.Set(src)
			}
			return nil
		}
	}
	return nil
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

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func DiscoverCri() (mlopsv1.CriType, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	nodeList, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{Limit: 1})
	if err != nil {
		return "", err
	}
	//nodeList := &v1.NodeList{}
	//if err := cr.List(context.Background(), nodeList, client.Limit(1)); err != nil {
	//	return "", err
	//}
	node := nodeList.Items[0]
	fmt.Println(node)
	return mlopsv1.CriTypeDocker, nil
}

func calculateAndApplyAppDefaults(app *mlopsv1.CnvrgApp, desiredAppSpec *mlopsv1.CnvrgAppSpec, infra *mlopsv1.CnvrgInfra) {
	if app.Spec.Cri == "" {
		DiscoverCri()
	}

	// set default heap size for ES if not set by user
	if strings.Contains(app.Spec.Dbs.Es.Requests.Memory, "Gi") && app.Spec.Dbs.Es.JavaOpts == "" {
		requestMem := strings.TrimSuffix(app.Spec.Dbs.Es.Requests.Memory, "Gi")
		mem, err := strconv.Atoi(requestMem)
		if err == nil {
			heapMem := mem / 2
			if heapMem > 0 {
				desiredAppSpec.Dbs.Es.JavaOpts = fmt.Sprintf("-Xms%dg -Xmx%dg", heapMem, heapMem)
			}
		}
	}

	if app.Spec.SSO.Enabled && app.Spec.ControlPlane.WebApp.Enabled {

		if app.Spec.ControlPlane.WebApp.OauthProxy.TokenValidationKey == "" {
			desiredAppSpec.ControlPlane.WebApp.OauthProxy.TokenValidationKey = generateSecureToken(16)
		}
		if app.Spec.ControlPlane.WebApp.OauthProxy.TokenValidationAuthData == "" {
			desiredAppSpec.ControlPlane.WebApp.OauthProxy.TokenValidationAuthData = generateSecureToken(6)
		}
	}

	if app.Spec.Networking.Ingress.IstioGwName == "" {
		desiredAppSpec.Networking.Ingress.IstioGwName = fmt.Sprintf(mlopsv1.IstioGwName, app.Namespace)
	}

	if app.Spec.Networking.Proxy.Enabled {
		desiredAppSpec.Networking.Proxy.NoProxy = app.Spec.Networking.Proxy.NoProxy
		// make sure no_proxy includes all default values
		for _, defaultNoProxy := range networking.DefaultNoProxy() {
			if !containsString(desiredAppSpec.Networking.Proxy.NoProxy, defaultNoProxy) {
				desiredAppSpec.Networking.Proxy.NoProxy = append(desiredAppSpec.Networking.Proxy.NoProxy, defaultNoProxy)
			}
		}
		// sort slices before compare
		sort.Strings(desiredAppSpec.Networking.Proxy.NoProxy)
		sort.Strings(app.Spec.Networking.Proxy.NoProxy)
		// if slice are not equal, use desiredAppSpec no_proxy
		if !reflect.DeepEqual(desiredAppSpec.Networking.Proxy.NoProxy, app.Spec.Networking.Proxy.NoProxy) {
			app.Spec.Networking.Proxy.NoProxy = nil
		}
	}

	if app.Spec.CnvrgAppPriorityClass.Name == "" && infra != nil {
		desiredAppSpec.CnvrgAppPriorityClass = infra.Spec.CnvrgAppPriorityClass
	}

	if app.Spec.CnvrgAppPriorityClass.Name == "" && infra != nil {
		desiredAppSpec.CnvrgJobPriorityClass = infra.Spec.CnvrgJobPriorityClass
	}
}

func calculateAndApplyInfraDefaults(infra *mlopsv1.CnvrgInfra, desiredInfraSpec *mlopsv1.CnvrgInfraSpec) {

	if infra.Spec.Networking.Ingress.IstioGwName == "" {
		desiredInfraSpec.Networking.Ingress.IstioGwName = fmt.Sprintf(mlopsv1.IstioGwName, infra.Spec.InfraNamespace)
	}

	if infra.Spec.Networking.Proxy.Enabled {
		desiredInfraSpec.Networking.Proxy.NoProxy = infra.Spec.Networking.Proxy.NoProxy
		for _, defaultNoProxy := range networking.DefaultNoProxy() {
			if !containsString(desiredInfraSpec.Networking.Proxy.NoProxy, defaultNoProxy) {
				desiredInfraSpec.Networking.Proxy.NoProxy = append(desiredInfraSpec.Networking.Proxy.NoProxy, defaultNoProxy)
			}
		}
		sort.Strings(desiredInfraSpec.Networking.Proxy.NoProxy)
		sort.Strings(infra.Spec.Networking.Proxy.NoProxy)
		if !reflect.DeepEqual(desiredInfraSpec.Networking.Proxy.NoProxy, infra.Spec.Networking.Proxy.NoProxy) {
			infra.Spec.Networking.Proxy.NoProxy = nil
		}
	}
}

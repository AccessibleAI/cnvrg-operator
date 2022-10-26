package app

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/controllers"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/networking"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/sso"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
	"strconv"
	"strings"
)

var log logr.Logger

func discoverOcpDefaultRouteHost(clientset client.Client) (ocpDefaultRouteHost string, err error) {
	routeCfg := &unstructured.Unstructured{}
	routeCfg.SetGroupVersionKind(desired.Kinds["OcpIngressCfgGVK"])
	routeCfg.SetName("cluster")
	err = clientset.Get(context.Background(), types.NamespacedName{Name: "cluster"}, routeCfg)
	if err != nil {
		return "", err
	}

	if _, ok := routeCfg.Object["spec"]; !ok {
		return "", fmt.Errorf("unable to parse OCP Ingress config, can't set default route host")
	}

	if domain, ok := routeCfg.Object["spec"].(map[string]interface{})["domain"]; !ok {
		return "", fmt.Errorf("unable to parse OCP Ingress config, can't set default route host")
	} else {
		return domain.(string), nil

	}
}

func CalculateAndApplyAppDefaults(app *mlopsv1.CnvrgApp, desiredAppSpec *mlopsv1.CnvrgAppSpec, clientset client.Client) error {

	if app.Spec.Dbs.Cvat.Enabled {
		desiredAppSpec.Dbs.Cvat.Pg.Enabled = true
		desiredAppSpec.Dbs.Cvat.Redis.Enabled = true
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

	if app.Spec.Networking.Ingress.IstioGwName == "" {
		desiredAppSpec.Networking.Ingress.IstioGwName = fmt.Sprintf("istio-gw-%s", app.Namespace)
	}

	if app.Spec.Networking.Proxy.Enabled {
		desiredAppSpec.Networking.Proxy.NoProxy = app.Spec.Networking.Proxy.NoProxy
		// make sure no_proxy includes all default values
		for _, defaultNoProxy := range networking.DefaultNoProxy(app.Spec.ClusterInternalDomain) {
			if !controllers.ContainsString(desiredAppSpec.Networking.Proxy.NoProxy, defaultNoProxy) {
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

	if app.Spec.Networking.Ingress.Type == mlopsv1.OpenShiftIngress {
		if app.Spec.ClusterDomain == "" && clientset != nil {
			clusterDomain, err := discoverOcpDefaultRouteHost(clientset)
			if err != nil {
				log.Error(err, "unable discover cluster domain, set clusterDomain manually under spec.clusterDomain")

			} else {
				desiredAppSpec.ClusterDomain = clusterDomain
			}
		}

	}

	if app.Spec.SSO.Enabled {
		if app.Spec.SSO.Central.PublicUrl == "" {
			scheme := "http"
			if app.Spec.Networking.HTTPS.Enabled {
				scheme = "https"
			}
			desiredAppSpec.SSO.Central.PublicUrl = fmt.Sprintf("%s://%s.%s", scheme, sso.CentralSsoSvcName, app.Spec.ClusterDomain)
		}
		if app.Spec.SSO.Proxy.Address == "" {
			desiredAppSpec.SSO.Proxy.Address = fmt.Sprintf("%s.%s.svc.%s",
				desiredAppSpec.SSO.Proxy.SvcName,
				app.Namespace,
				desiredAppSpec.ClusterInternalDomain,
			)
		}
		if app.Spec.SSO.Authz.Address == "" {
			desiredAppSpec.SSO.Authz.Address = fmt.Sprintf("%s.%s.svc:50052", desiredAppSpec.SSO.Authz.SvcName, app.Namespace)
		}
	}

	return nil
}

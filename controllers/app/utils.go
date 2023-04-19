package app

import (
	"context"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/controllers"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/networking"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
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

func CalculateAndApplyAppDefaults(app *mlopsv1.CnvrgApp, defaultSpec *mlopsv1.CnvrgAppSpec, clientset client.Client) error {

	if app.Spec.Dbs.Cvat.Enabled {
		defaultSpec.Dbs.Cvat.Pg.Enabled = true
		defaultSpec.Dbs.Cvat.Redis.Enabled = true
	}

	// set default heap size for ES if not set by user
	if strings.Contains(app.Spec.Dbs.Es.Requests.Memory, "Gi") && app.Spec.Dbs.Es.JavaOpts == "" {
		requestMem := strings.TrimSuffix(app.Spec.Dbs.Es.Requests.Memory, "Gi")
		mem, err := strconv.Atoi(requestMem)
		if err == nil {
			heapMem := mem / 2
			if heapMem > 0 {
				defaultSpec.Dbs.Es.JavaOpts = fmt.Sprintf("-Xms%dg -Xmx%dg", heapMem, heapMem)
			}
		}
	}

	if app.Spec.Networking.Ingress.IstioGwName == "" {
		defaultSpec.Networking.Ingress.IstioGwName = fmt.Sprintf("istio-gw-%s", app.Namespace)
	}

	if app.Spec.Networking.Proxy.Enabled {
		defaultSpec.Networking.Proxy.NoProxy = app.Spec.Networking.Proxy.NoProxy
		// make sure no_proxy includes all default values
		for _, defaultNoProxy := range networking.DefaultNoProxy(app.Spec.ClusterInternalDomain) {
			if !controllers.ContainsString(defaultSpec.Networking.Proxy.NoProxy, defaultNoProxy) {
				defaultSpec.Networking.Proxy.NoProxy = append(defaultSpec.Networking.Proxy.NoProxy, defaultNoProxy)
			}
		}
		// sort slices before compare
		sort.Strings(defaultSpec.Networking.Proxy.NoProxy)
		sort.Strings(app.Spec.Networking.Proxy.NoProxy)
		// if slice are not equal, use defaultSpec no_proxy
		if !reflect.DeepEqual(defaultSpec.Networking.Proxy.NoProxy, app.Spec.Networking.Proxy.NoProxy) {
			app.Spec.Networking.Proxy.NoProxy = nil
		}
	}

	if app.Spec.Networking.Ingress.Type == mlopsv1.OpenShiftIngress {
		if app.Spec.ClusterDomain == "" && clientset != nil {
			clusterDomain, err := discoverOcpDefaultRouteHost(clientset)
			if err != nil {
				log.Error(err, "unable discover cluster domain, set clusterDomain manually under spec.clusterDomain")

			} else {
				defaultSpec.ClusterDomain = clusterDomain
				app.Spec.ClusterDomain = clusterDomain
			}
		}

	}

	if app.Spec.SSO.Enabled {
		scheme := "http"
		if app.Spec.Networking.HTTPS.Enabled {
			scheme = "https"
		}
		if app.Spec.SSO.Central.PublicUrl == "" {
			defaultSpec.SSO.Central.PublicUrl = fmt.Sprintf("%s://%s.%s", scheme, defaultSpec.SSO.Central.SvcName, app.Spec.ClusterDomain)
		}
		if app.Spec.SSO.Proxy.Address == "" {
			defaultSpec.SSO.Proxy.Address = fmt.Sprintf("%s.%s.svc.%s",
				defaultSpec.SSO.Proxy.SvcName,
				app.Namespace,
				defaultSpec.ClusterInternalDomain,
			)
		}
		if app.Spec.SSO.Central.JwksURL == "" {

			defaultSpec.SSO.Central.JwksURL = fmt.Sprintf("%s://%s.%s/v1/%s/.well-known/jwks.json?client_id", scheme,
				defaultSpec.SSO.Jwks.SvcName,
				app.Spec.ClusterDomain,
				strings.Split(app.Spec.ClusterDomain, ".")[0],
			)

		} else if !strings.HasSuffix(app.Spec.SSO.Central.JwksURL, "jwks.json?client_id") {
			app.Spec.SSO.Central.JwksURL = fmt.Sprintf("%s/v1/%s/.well-known/jwks.json?client_id",
				app.Spec.SSO.Central.JwksURL,
				strings.Split(app.Spec.ClusterDomain, ".")[0],
			)
		}
	}

	return nil
}

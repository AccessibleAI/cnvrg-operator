package app

import (
	"context"
	"encoding/json"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/controllers"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/networking"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func discoverOcpPromCreds(clientset client.Client, ns string) {
	promCreds := &v1.Secret{ObjectMeta: metav1.ObjectMeta{
		Name:      "prom-creds",
		Namespace: ns,
	}}
	err := clientset.Get(context.Background(), types.NamespacedName{Name: promCreds.Name, Namespace: promCreds.Namespace}, promCreds)
	if errors.IsNotFound(err) {
		dataSourceSecret := &v1.Secret{}
		name := types.NamespacedName{Name: "grafana-datasources-v2", Namespace: "openshift-monitoring"}
		err := clientset.Get(context.Background(), name, dataSourceSecret)
		if err != nil {
			log.Error(err, "failed to discover prom creds")
			return

		}
		if _, ok := dataSourceSecret.Data["prometheus.yaml"]; !ok {
			log.Error(err, "failed to discover prom creds")
			return
		}
		promData := &struct {
			Datasources []struct {
				Url            string
				BasicAuthUser  string
				SecureJsonData struct {
					BasicAuthPassword string
				}
			}
		}{}

		err = json.Unmarshal(dataSourceSecret.Data["prometheus.yaml"], &promData)

		if err != nil {
			log.Error(err, "error unmarshal prometheus.yaml")
			return
		}

		if len(promData.Datasources) < 1 {
			log.Error(err, "unexpected prom data")
			return
		}

		promCreds.Data = make(map[string][]byte)
		promCreds.Data["CNVRG_PROMETHEUS_URL"] = []byte(promData.Datasources[0].Url)
		promCreds.Data["CNVRG_PROMETHEUS_USER"] = []byte(promData.Datasources[0].BasicAuthUser)
		promCreds.Data["CNVRG_PROMETHEUS_PASS"] = []byte(promData.Datasources[0].SecureJsonData.BasicAuthPassword)

		err = clientset.Create(context.Background(), promCreds)
		if err != nil {
			log.Error(err, "error creating promCreds for OCP setup")
		}

	} else if err != nil {
		log.Error(err, "failed to discover prom creds")
	}
	return
}

func calculateAndApplyAppDefaults(app *mlopsv1.CnvrgApp, desiredAppSpec *mlopsv1.CnvrgAppSpec, clientset client.Client) error {

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
		if app.Spec.ClusterDomain == "" {
			clusterDomain, err := discoverOcpDefaultRouteHost(clientset)
			if err != nil {
				log.Error(err, "unable discover cluster domain, set clusterDomain manually under spec.clusterDomain")

			} else {
				desiredAppSpec.ClusterDomain = clusterDomain
			}
		}
		if !app.Spec.Dbs.Prom.Enabled {
			discoverOcpPromCreds(clientset, app.Namespace)
		}
	}

	return nil
}

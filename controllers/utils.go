package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/networking"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
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

func discoverCri(clientset client.Client) (mlopsv1.CriType, error) {
	nodeList := &v1.NodeList{}
	if err := clientset.List(context.Background(), nodeList, client.Limit(1)); err != nil {
		return "", err
	}
	if len(nodeList.Items) == 0 {
		return "", errors.New("could not recognize cri because we could not list nodes. you can set it manually in cnvrgapp and cnvrginfra")
	}
	node := nodeList.Items[0]
	cri := node.Status.NodeInfo.ContainerRuntimeVersion
	if strings.Contains(cri, string(mlopsv1.CriTypeContainerd)) {
		return mlopsv1.CriTypeContainerd, nil
	} else if strings.Contains(cri, string(mlopsv1.CriTypeCrio)) {
		return mlopsv1.CriTypeCrio, nil
	} else if strings.Contains(cri, string(mlopsv1.CriTypeDocker)) {
		return mlopsv1.CriTypeDocker, nil
	} else {
		return "", errors.New("could not recognize cri. you can set it manually in cnvrgapp and cnvrginfra")
	}
}

func discoverOcpDefaultRouteHost(clientset client.Client) (ocpDefaultRouteHost string, err error) {
	routeCfg := &unstructured.Unstructured{}
	routeCfg.SetGroupVersionKind(desired.Kinds["OcpIngressCfg"])
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
	if k8serrors.IsNotFound(err) {
		dataSourceSecret := &v1.Secret{}
		name := types.NamespacedName{Name: "grafana-datasources-v2", Namespace: "openshift-monitoring"}
		err := clientset.Get(context.Background(), name, dataSourceSecret)
		if err != nil {
			appLog.Error(err, "failed to discover prom creds")
			return

		}
		if _, ok := dataSourceSecret.Data["prometheus.yaml"]; !ok {
			appLog.Error(err, "failed to discover prom creds")
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
			appLog.Error(err, "error unmarshal prometheus.yaml")
			return
		}

		if len(promData.Datasources) < 1 {
			appLog.Error(err, "unexpected prom data")
			return
		}

		promCreds.Data = make(map[string][]byte)
		promCreds.Data["CNVRG_PROMETHEUS_URL"] = []byte(promData.Datasources[0].Url)
		promCreds.Data["CNVRG_PROMETHEUS_USER"] = []byte(promData.Datasources[0].BasicAuthUser)
		promCreds.Data["CNVRG_PROMETHEUS_PASS"] = []byte(promData.Datasources[0].SecureJsonData.BasicAuthPassword)

		err = clientset.Create(context.Background(), promCreds)
		if err != nil {
			appLog.Error(err, "error creating promCreds for OCP setup")
		}

	} else if err != nil {
		appLog.Error(err, "failed to discover prom creds")
	}
	return
}

func calculateAndApplyAppDefaults(app *mlopsv1.CnvrgApp, desiredAppSpec *mlopsv1.CnvrgAppSpec, infra *mlopsv1.CnvrgInfra, clientset client.Client) error {
	if app.Spec.Cri == "" {
		cri, err := discoverCri(clientset)
		if err != nil {
			return err
		}
		desiredAppSpec.Cri = cri
	}

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
		desiredAppSpec.Networking.Ingress.IstioGwName = fmt.Sprintf(mlopsv1.IstioGwName, app.Namespace)
	}

	if app.Spec.Networking.Proxy.Enabled {
		desiredAppSpec.Networking.Proxy.NoProxy = app.Spec.Networking.Proxy.NoProxy
		// make sure no_proxy includes all default values
		for _, defaultNoProxy := range networking.DefaultNoProxy(app.Spec.ClusterInternalDomain) {
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

	if app.Spec.Networking.Ingress.Type == mlopsv1.OpenShiftIngress {
		if app.Spec.ClusterDomain == "" {
			clusterDomain, err := discoverOcpDefaultRouteHost(clientset)
			if err != nil {
				appLog.Error(err, "unable discover cluster domain, set clusterDomain manually under spec.clusterDomain")

			} else {
				desiredAppSpec.ClusterDomain = clusterDomain
			}
		}
		discoverOcpPromCreds(clientset, app.Namespace)
	}

	return nil
}

func calculateAndApplyInfraDefaults(infra *mlopsv1.CnvrgInfra, desiredInfraSpec *mlopsv1.CnvrgInfraSpec, clientset client.Client) error {
	if infra.Spec.Cri == "" {
		cri, err := discoverCri(clientset)
		if err != nil {
			return err
		}

		desiredInfraSpec.Cri = cri
	}

	if infra.Spec.Networking.Ingress.IstioGwName == "" {
		desiredInfraSpec.Networking.Ingress.IstioGwName = fmt.Sprintf(mlopsv1.IstioGwName, infra.Spec.InfraNamespace)
	}

	if infra.Spec.Networking.Proxy.Enabled {
		desiredInfraSpec.Networking.Proxy.NoProxy = infra.Spec.Networking.Proxy.NoProxy
		for _, defaultNoProxy := range networking.DefaultNoProxy(infra.Spec.ClusterInternalDomain) {
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

	return nil
}

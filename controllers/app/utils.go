package app

import (
	"context"
	"crypto/sha256"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/controllers"
	"github.com/AccessibleAI/cnvrg-operator/pkg/app/networking"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/go-logr/logr"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	mathrand "math/rand"
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

func getTlsCert(app *mlopsv1.CnvrgApp, clientset client.Client) (cert string, key string, err error) {
	namespacedName := types.NamespacedName{Name: app.Spec.Networking.HTTPS.CertSecret, Namespace: app.Namespace}
	certSecret := v1core.Secret{ObjectMeta: metav1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := clientset.Get(context.Background(), namespacedName, &certSecret); err != nil {
		return "", "", err
	}

	if _, ok := certSecret.Data["tls.crt"]; !ok {
		err := fmt.Errorf("certificate secret %s missing required field tls.crt", namespacedName.Name)
		return "", "", err
	}

	if _, ok := certSecret.Data["tls.key"]; !ok {
		err := fmt.Errorf("certificate secret %s missing required field tls.key", namespacedName.Name)
		return "", "", err
	}

	return string(certSecret.Data["tls.crt"]), string(certSecret.Data["tls.key"]), nil
}

func CalculateAndApplyAppDefaults(app *mlopsv1.CnvrgApp, defaultSpec *mlopsv1.CnvrgAppSpec, clientset client.Client) error {

	// set cluster domain prefix if enabled
	setClusterDomainPrefix(app, defaultSpec)

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

	// discover defaults for OpenShift Route Ingress
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

	if app.Spec.Networking.Ingress.OcpSecureRoutes &&
		(app.Spec.Networking.HTTPS.Cert == "" || app.Spec.Networking.HTTPS.Key == "" ||
			defaultSpec.Networking.HTTPS.CertSecret != app.Spec.Networking.HTTPS.CertSecret) {
		cert, key, err := getTlsCert(app, clientset)
		if err != nil {
			log.Error(err, "unable to retrieve tls secret")
		} else {
			defaultSpec.Networking.HTTPS.Cert = cert
			defaultSpec.Networking.HTTPS.Key = key
		}
	}

	// configure defaults for SSO
	if app.Spec.SSO.Enabled {
		scheme := "http"
		if app.Spec.Networking.HTTPS.Enabled {
			scheme = "https"
		}
		if app.Spec.SSO.Central.PublicUrl == "" {
			defaultSpec.SSO.Central.PublicUrl = fmt.Sprintf("%s://%s%s.%s",
				scheme,
				defaultSpec.SSO.Central.SvcName,
				app.Spec.Networking.ClusterDomainPrefix.Prefix,
				app.Spec.ClusterDomain)
		}
		if app.Spec.SSO.Proxy.Address == "" {
			defaultSpec.SSO.Proxy.Address = fmt.Sprintf("%s.%s.svc.%s",
				defaultSpec.SSO.Proxy.SvcName,
				app.Namespace,
				defaultSpec.ClusterInternalDomain,
			)
		}
		if app.Spec.SSO.Central.JwksURL == "" {

			defaultSpec.SSO.Central.JwksURL = fmt.Sprintf("%s://%s%s.%s/v1/%s/.well-known/jwks.json?client_id", scheme,
				defaultSpec.SSO.Jwks.SvcName,
				app.Spec.Networking.ClusterDomainPrefix.Prefix,
				app.Spec.ClusterDomain,
				strings.Split(app.Spec.ClusterDomain, ".")[0],
			)

		} else if !strings.HasSuffix(app.Spec.SSO.Central.JwksURL, "jwks.json?client_id") {
			app.Spec.SSO.Central.JwksURL = fmt.Sprintf("%s/v1/%s/.well-known/jwks.json?client_id",
				app.Spec.SSO.Central.JwksURL,
				strings.Split(app.Spec.ClusterDomain, ".")[0],
			)
		}
		if app.Spec.SSO.Central.WhitelistDomain == "" {
			defaultSpec.SSO.Central.WhitelistDomain = fmt.Sprintf(".%s", app.Spec.ClusterDomain)
		}
		if app.Spec.SSO.Central.CookieDomain == "" {
			defaultSpec.SSO.Central.CookieDomain = fmt.Sprintf(".%s", app.Spec.ClusterDomain)
		}
	}

	return nil
}

func setClusterDomainPrefix(app *mlopsv1.CnvrgApp, defaultSpec *mlopsv1.CnvrgAppSpec) {
	if app.Spec.Networking.ClusterDomainPrefix.Enabled && app.Spec.Networking.ClusterDomainPrefix.Prefix == "" {
		defaultSpec.Networking.ClusterDomainPrefix.Prefix = fmt.Sprintf("-%s", RandomString(5))
	}
}

func labelsMapToList(labels map[string]string) (labelList []string) {
	for labelName, _ := range labels {
		labelList = append(labelList, labelName)
	}
	return labelList
}

func RandomString(length int) string {
	var output strings.Builder
	charSet := "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < length; i++ {
		random := mathrand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

func hashStringsMap(m map[string]string) string {
	h := sha256.New()

	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]

		b := sha256.Sum256([]byte(fmt.Sprintf("%v", k)))
		h.Write(b[:])
		b = sha256.Sum256([]byte(fmt.Sprintf("%v", v)))
		h.Write(b[:])
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

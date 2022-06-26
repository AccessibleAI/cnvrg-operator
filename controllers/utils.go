package controllers

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/pkg/networking"
	v1 "k8s.io/api/core/v1"
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

func generateKeys() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	// private key
	privatePemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	privatePemBytes := pem.EncodeToMemory(privatePemBlock)

	// public key
	b, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	publicPemBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: b}
	publicPemBytes := pem.EncodeToMemory(publicPemBlock)

	return privatePemBytes, publicPemBytes, nil
}

func DiscoverCri(clientset client.Client) (mlopsv1.CriType, error) {
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

func calculateAndApplyAppDefaults(app *mlopsv1.CnvrgApp, desiredAppSpec *mlopsv1.CnvrgAppSpec, infra *mlopsv1.CnvrgInfra, clientset client.Client) error {
	if app.Spec.Cri == "" {
		cri, err := DiscoverCri(clientset)
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

	if app.Spec.SSO.Enabled {
		if desiredAppSpec.ControlPlane.BaseConfig.FeatureFlags == nil {
			desiredAppSpec.ControlPlane.BaseConfig.FeatureFlags = make(map[string]string)
		}
		desiredAppSpec.ControlPlane.BaseConfig.FeatureFlags["JWKS_ISS"] = getJWTIss(app, infra)
		desiredAppSpec.ControlPlane.BaseConfig.FeatureFlags["JWKS_AUD"] = getJWTAud()
		desiredAppSpec.SSO.ExtraJWTIssuers = append(app.Spec.SSO.ExtraJWTIssuers, getExtraJWTIssuers(app, infra)...)

		if app.Spec.SSO.SaaSSSO.Enabled {
			desiredAppSpec.SSO.SaaSSSO.AllowedGroups = append(app.Spec.SSO.SaaSSSO.AllowedGroups, getTenantGroup(app.Spec.ClusterDomain)...)
		}

	}

	return nil
}

func getTenantGroup(clusterDomain string) []string {
	var tenantGroups []string
	group := strings.Split(clusterDomain, ".")[0]
	// regular group
	tenantGroups = append(tenantGroups, group)
	// Keycloak group - always prefixed with a slash '/`
	tenantGroups = append(tenantGroups, fmt.Sprintf(`/%s`, group))
	return tenantGroups
}

func getJWTIss(app *mlopsv1.CnvrgApp, infra *mlopsv1.CnvrgInfra) string {
	jwksSvcName := "cnvrg-jwks"
	if infra.Spec.Jwks.Name != "" {
		jwksSvcName = infra.Spec.Jwks.Name
	}
	issuerAddress := fmt.Sprintf("%s.%s/v1/%s/.well-known/jwks.json?client_id",
		jwksSvcName,
		infra.Spec.ClusterDomain,
		strings.Split(app.Spec.ClusterDomain, ".")[0],
	)
	if app.Spec.Networking.HTTPS.Enabled {
		return fmt.Sprintf("https://%s", issuerAddress)
	} else {
		return fmt.Sprintf("http://%s", issuerAddress)
	}
}

func getJWTAud() string {
	return "cnvrg-tenant"
}
func getExtraJWTIssuers(app *mlopsv1.CnvrgApp, infra *mlopsv1.CnvrgInfra) (issuers []string) {
	issuers = append(issuers, fmt.Sprintf("%s=%s", getJWTIss(app, infra), getJWTAud()))
	return
}

func calculateAndApplyInfraDefaults(infra *mlopsv1.CnvrgInfra, desiredInfraSpec *mlopsv1.CnvrgInfraSpec, clientset client.Client) error {
	if infra.Spec.Cri == "" {
		cri, err := DiscoverCri(clientset)
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

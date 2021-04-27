package desired

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Dimss/crypt/apr1_crypt"
	"github.com/Masterminds/sprig"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	mathrand "math/rand"
	"os"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var GrafanaAppDashboards = []string{
	"k8s-resources-namespace.json",
	"k8s-resources-pod.json",
	"k8s-resources-workload.json",
	"k8s-resources-workloads-namespace.json",
	"namespace-by-pod.json",
	"namespace-by-workload.json",
	"persistentvolumesusage.json",
	"pod-total.json",
	"statefulset.json",
	"workload-total.json",
	"grafana-idle-metrics.json",
}

var GrafanaInfraDashboards = append([]string{
	"apiserver.json",
	"cluster-total.json",
	"controller-manager.json",
	"k8s-resources-cluster.json",
	"k8s-resources-node.json",
	"kubelet.json",
	"node-cluster-rsrc-use.json",
	"node-rsrc-use.json",
	"nodes.json",
	"prometheus-remote-write.json",
	"prometheus.json",
	"proxy.json",
	"scheduler.json",
	"node-exporter.json",
}, GrafanaAppDashboards...)

func getNs(obj interface{}) string {
	if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(obj) {
		cnvrgInfra := obj.(*mlopsv1.CnvrgInfra)
		return cnvrgInfra.Spec.InfraNamespace
	}
	if reflect.TypeOf(&mlopsv1.CnvrgApp{}) == reflect.TypeOf(obj) {
		cnvrgApp := obj.(*mlopsv1.CnvrgApp)
		return cnvrgApp.Namespace
	}
	return "cnvrg-infra"
}

func getIstioGwName(obj interface{}) string {
	if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(obj) {
		cnvrgInfra := obj.(*mlopsv1.CnvrgInfra)
		if cnvrgInfra.Spec.Networking.Ingress.IstioGwName == "" {
			return fmt.Sprintf("isito-gw-%v", getNs(obj))
		} else {
			return cnvrgInfra.Spec.Networking.Ingress.IstioGwName
		}
	}
	if reflect.TypeOf(&mlopsv1.CnvrgApp{}) == reflect.TypeOf(obj) {
		cnvrgApp := obj.(*mlopsv1.CnvrgApp)
		if cnvrgApp.Spec.Networking.Ingress.IstioGwName == "" {
			return fmt.Sprintf("isito-gw-%v", getNs(obj))
		} else {
			return cnvrgApp.Spec.Networking.Ingress.IstioGwName
		}
	}
	return "" // what can go wrong? :)
}

func getGrafanaDashboards(obj interface{}) []string {
	if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(obj) {
		return GrafanaInfraDashboards
	}
	if reflect.TypeOf(&mlopsv1.CnvrgApp{}) == reflect.TypeOf(obj) {
		cnvrgApp := obj.(*mlopsv1.CnvrgApp)
		if cnvrgApp.Spec.NamespaceTenancy == false {
			return GrafanaInfraDashboards
		}
		return GrafanaAppDashboards
	}
	return nil
}

func getSSOConfig(obj interface{}) *mlopsv1.SSO {
	if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(obj) {
		return &obj.(*mlopsv1.CnvrgInfra).Spec.SSO
	}
	if reflect.TypeOf(&mlopsv1.CnvrgApp{}) == reflect.TypeOf(obj) {
		return &obj.(*mlopsv1.CnvrgApp).Spec.SSO
	}
	return nil
}

func getSSORedirectUrl(obj interface{}, svc string) string {
	if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(obj) {
		infra := obj.(*mlopsv1.CnvrgInfra)
		if *infra.Spec.Networking.HTTPS.Enabled {
			return fmt.Sprintf("https://%v.%v/oauth2/callback", svc, infra.Spec.ClusterDomain)
		} else {
			return fmt.Sprintf("http://%v.%v/oauth2/callback", svc, infra.Spec.ClusterDomain)
		}

	}
	if reflect.TypeOf(&mlopsv1.CnvrgApp{}) == reflect.TypeOf(obj) {
		app := obj.(*mlopsv1.CnvrgApp)
		if *app.Spec.Networking.HTTPS.Enabled {
			return fmt.Sprintf("https://%v.%v/oauth2/callback", svc, app.Spec.ClusterDomain)
		} else {
			return fmt.Sprintf("http://%v.%v/oauth2/callback", svc, app.Spec.ClusterDomain)
		}
	}
	return ""
}

func cnvrgTemplateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"ns": func(obj interface{}) string {
			return getNs(obj)
		},
		"httpScheme": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if *cnvrgApp.Spec.Networking.HTTPS.Enabled {
				return "https://"
			}
			return "http://"
		},
		"appDomain": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.NodePortIngress {
				return cnvrgApp.Spec.ClusterDomain + ":" +
					strconv.Itoa(cnvrgApp.Spec.ControlPlane.WebApp.NodePort)
			} else {
				return cnvrgApp.Spec.ControlPlane.WebApp.SvcName + "." + cnvrgApp.Spec.ClusterDomain
			}
		},
		"defaultComputeClusterDomain": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.NodePortIngress {
				return cnvrgApp.Spec.ClusterDomain + ":" +
					strconv.Itoa(cnvrgApp.Spec.ControlPlane.WebApp.NodePort)
			} else {
				return cnvrgApp.Spec.ClusterDomain
			}
		},
		"redisUrl": func(cnvrgApp mlopsv1.CnvrgApp) string {
			return "redis://" + cnvrgApp.Spec.Dbs.Redis.SvcName
		},
		"esUrl": func(cnvrgApp mlopsv1.CnvrgApp) string {
			return "http://" + cnvrgApp.Spec.Dbs.Es.SvcName
		},
		"hyperServerUrl": func(cnvrgApp mlopsv1.CnvrgApp) string {
			return "http://" + cnvrgApp.Spec.ControlPlane.Hyper.SvcName
		},
		"esFullInternalUrl": func(cnvrgApp mlopsv1.CnvrgApp) string {
			return fmt.Sprintf("http://%s.%s.svc.cluster.local:%d",
				cnvrgApp.Spec.Dbs.Es.SvcName,
				cnvrgApp.Namespace,
				cnvrgApp.Spec.Dbs.Es.Port)
		},
		"objectStorageUrl": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.ControlPlane.ObjectStorage.CnvrgStorageEndpoint != "" {
				return cnvrgApp.Spec.ControlPlane.ObjectStorage.CnvrgStorageEndpoint
			}
			if *cnvrgApp.Spec.Networking.HTTPS.Enabled {
				return fmt.Sprintf("https://%s.%s", cnvrgApp.Spec.Dbs.Minio.SvcName, cnvrgApp.Spec.ClusterDomain)
			} else {
				return fmt.Sprintf("http://%s.%s", cnvrgApp.Spec.Dbs.Minio.SvcName, cnvrgApp.Spec.ClusterDomain)
			}
		},
		"routeBy": func(cnvrgApp mlopsv1.CnvrgApp, routeBy string) string {
			switch routeBy {
			case "ISTIO":
				if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.IstioIngress {
					return "true"
				}
				return "false"
			case "OPENSHIFT":
				if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.OpenShiftIngress {
					return "true"
				}
				return "false"
			case "NGINX_INGRESS":
				if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.NginxIngress {
					return "true"
				}
				return "false"
			case "NODE_PORT":
				if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.NodePortIngress {
					return "true"
				}
				return "false"
			}
			return "false"
		},
		"oauthProxyConfig": func(obj interface{}, svc string, skipAuthRegex []string, provider string, proxyPort, upstreamPort int) string {
			sso := getSSOConfig(obj)
			skipAuthUrls := fmt.Sprintf(`["%v", `, `^\/cnvrg-static/`)
			for i, url := range skipAuthRegex {
				if i == (len(skipAuthRegex) - 1) {
					skipAuthUrls += fmt.Sprintf(`"%v"`, url)
				} else {
					skipAuthUrls += fmt.Sprintf(`"%v", `, url)
				}
			}
			skipAuthUrls += "]"
			proxyConf := []string{
				fmt.Sprintf(`provider = "%v"`, provider),
				fmt.Sprintf(`http_address = "0.0.0.0:%d"`, proxyPort),
				fmt.Sprintf(`redirect_url = "%v"`, getSSORedirectUrl(obj, svc)),
				fmt.Sprintf("skip_auth_regex = %v", skipAuthUrls),
				fmt.Sprintf(`email_domains = ["%v"]`, sso.EmailDomain),
				fmt.Sprintf(`client_id = "%v"`, sso.ClientID),
				fmt.Sprintf(`client_secret = "%v"`, sso.ClientSecret),
				fmt.Sprintf(`cookie_secret = "%v"`, sso.CookieSecret),
				fmt.Sprintf(`oidc_issuer_url = "%v"`, sso.OidcIssuerURL),
				fmt.Sprintf(`upstreams = ["http://127.0.0.1:%d/", "file:///opt/app-root/src/templates/#/cnvrg-static/"]`, upstreamPort),
				`session_store_type = "redis"`,
				`skip_jwt_bearer_tokens = true`,
				`custom_templates_dir = "/opt/app-root/src/templates"`,
				"ssl_insecure_skip_verify = true",
				`cookie_name = "_oauth2_proxy"`,
				`cookie_expire = "168h"`,
				"cookie_secure = false",
				"cookie_httponly = true",
			}

			return strings.Join(proxyConf, "\n")
		},
		"cnvrgPassengerBindAddress": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if *cnvrgApp.Spec.SSO.Enabled {
				return "127.0.0.1"
			}
			return "0.0.0.0"
		},
		"cnvrgPassengerBindPort": func(cnvrgApp mlopsv1.CnvrgApp) int {
			if *cnvrgApp.Spec.SSO.Enabled {
				return 3000
			}
			return cnvrgApp.Spec.ControlPlane.WebApp.Port
		},
		"prometheusStaticConfig": func(cnvrgApp mlopsv1.CnvrgApp, ns string) string {
			if cnvrgApp.Spec.NamespaceTenancy == true {
				return fmt.Sprintf(`
- job_name: 'federate'
  scrape_interval: 10s
  honor_labels: true
  honor_timestamps: false
  metrics_path: '/federate'
  basic_auth:
    username: 'email@username.me'
    password: 'cfgqvzjbhnwcomplicatedpasswordwjnqmd'
  params:
    'match[]':
      - '{namespace="%s"}'
  static_configs:
    - targets:
      - '%s'
`, ns, "asd")
			}
			return fmt.Sprintf(`
- job_name: 'federate'
  scrape_interval: 10s
  honor_labels: true
  honor_timestamps: false
  metrics_path: '/federate'
  params:
    'match[]':
      - '{job=~".+"}'
  static_configs:
    - targets:
      - '%s'
`, "cnvrgApp.Spec.Monitoring.UpstreamPrometheus")
		},
		"grafanaDataSource": func(promUrl, user, pass string) string {
			return fmt.Sprintf(`
{
    "apiVersion": 1,
    "datasources": [
        {
            "access": "proxy",
            "editable": false,
            "name": "prometheus",
            "orgId": 1,
            "type": "prometheus",
            "url": "%s",
            "version": 1,
            "basicAuth": true,
            "basicAuthUser": "%s",
            "basicAuthPassword": "%s",
            "secureJsonFields": {
              "basicAuthPassword": true
            },
            "jsonData": {
                "tlsSkipVerify": true
            },
        }
    ]
}`, promUrl, user, pass)
		},
		"grafanaDashboards": func(obj interface{}) []string {
			return getGrafanaDashboards(obj)
		},
		"isAppSpec": func(obj interface{}) bool {
			if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(obj) {
				return false
			}
			return true
		},
		"istioGwName": func(obj interface{}) string {
			return getIstioGwName(obj)
		},
		"kibanaSecret": func(host, port, esHost, esUser, esPass, esBasicAuth string) string {
			return fmt.Sprintf(`
server:
  name: kibana
  host: %s
  port: %s
elasticsearch:
  hosts:
  - %s
  username: %s
  password: %s
  customHeaders:
    Authorization: "Basic %s"
`, host, port, esHost, esUser, esPass, esBasicAuth)
		},
		"isTrue": func(boolPointer *bool) bool {
			return *boolPointer
		},
	}
}

func (s *State) GenerateDeployable() error {
	var tpl bytes.Buffer
	f, err := pkger.Open(s.TemplatePath)
	if err != nil {
		zap.S().Error(err, "error reading path", "path", s.TemplatePath)
		return err
	}
	b, err := ioutil.ReadAll(f)

	if err != nil {
		zap.S().Errorf("%v, error reading file: %v", err, s.TemplatePath)
		return err
	}

	s.Template, err = template.New(strings.ReplaceAll(s.TemplatePath, "/", "-")).
		Funcs(sprig.TxtFuncMap()).
		Funcs(cnvrgTemplateFuncs()).
		Parse(string(b))
	if err != nil {
		zap.S().Errorf("%v, template: %v", err, s.TemplatePath)
		return err
	}
	s.Obj.SetGroupVersionKind(s.GVR)
	if err := s.Template.Execute(&tpl, s.TemplateData); err != nil {
		zap.S().Error(err, "rendering template error", "file", s.TemplatePath)
		return err
	}
	s.ParsedTemplate = tpl.String()
	zap.S().Debug("parsing: %v ", s.TemplatePath)
	zap.S().Debug("template: " + s.TemplatePath + "\n" + s.ParsedTemplate)
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	if _, _, err := dec.Decode([]byte(s.ParsedTemplate), nil, s.Obj); err != nil {
		zap.S().Errorf("%v, template: %v", err, s.ParsedTemplate)
		return err
	}
	if err := s.dumpTemplateToFile(); err != nil {
		zap.S().Error(err, "dumping template file", "file", s.TemplatePath)
		return err
	}
	return nil
}

func Apply(desiredManifests []*State, desiredSpec v1.Object, client client.Client, schema *runtime.Scheme, log logr.Logger) error {

	ctx := context.Background()
	for _, manifest := range desiredManifests {

		if manifest.TemplateData == nil {

			manifest.TemplateData = desiredSpec
		}
		if err := manifest.GenerateDeployable(); err != nil {
			log.Error(err, "error generating deployable", "name", manifest.Obj.GetName())
			return err
		}

		if manifest.Own {
			if err := ctrl.SetControllerReference(desiredSpec, manifest.Obj, schema); err != nil {
				log.Error(err, "error setting controller reference", "name", manifest.Obj.GetName())
				return err
			}
		}
		if viper.GetBool("dry-run") {
			log.Info("dry run enabled, skipping applying...")
			continue
		}
		fetchInto := &unstructured.Unstructured{}
		fetchInto.SetGroupVersionKind(manifest.GVR)
		err := client.Get(ctx, types.NamespacedName{Name: manifest.Obj.GetName(), Namespace: manifest.Obj.GetNamespace()}, fetchInto)
		if err != nil && errors.IsNotFound(err) {
			log.V(1).Info("creating", "name", manifest.Obj.GetName(), "kind", manifest.GVR.Kind)
			if err := client.Create(ctx, manifest.Obj); err != nil {
				log.Error(err, "error creating object", "name", manifest.Obj.GetName())
				return err
			}
		} else {

			if !shouldUpdate(manifest, fetchInto) {
				continue
			}

			if err := mergo.Merge(fetchInto, manifest.Obj, mergo.WithOverride); err != nil {
				log.Error(err, "can't merge")
				return err
			}

			finalObjToApply := fetchInto

			if manifest.Override {
				finalObjToApply = manifest.Obj // if override true, do not merge object with existing state
			}

			err := client.Update(ctx, finalObjToApply)
			if err != nil {
				log.Info("error updating object", "manifest", manifest.TemplatePath)
				return err
			}
		}
	}
	return nil
}

func shouldUpdate(manifest *State, obj *unstructured.Unstructured) bool {

	// do not try to update PVC, they are immutable (probably)
	if manifest.GVR == Kinds[PvcGVR] {
		return false
	}

	// todo: figure out what to do with this (happens only on OCP, why?)
	if manifest.GVR == Kinds[SaGVR] {
		return false
	}

	// do not apply CRDs if already exists
	// todo: have to ensure that existing CRD version is compatible with actually CR
	if manifest.GVR == Kinds[CrdGVR] {
		return false
	}

	// todo: need to figure out what's wrong with MPI CRD (might be related to apiextensions.k8s.io/v1beta1)
	if manifest.GVR == Kinds[CrdGVR] && obj.GetName() == "mpijobs.kubeflow.org" {
		return false
	}

	return true
}

func (s *State) dumpTemplateToFile() error {
	templatesDumpDir := viper.GetString("templates-dump-dir")
	if templatesDumpDir != "" {
		if _, err := os.Stat(templatesDumpDir); os.IsNotExist(err) {
			if err = os.Mkdir(templatesDumpDir, 0775); err != nil {
				zap.S().Error(err, "can't create templates dump dir for templates debugging")
				return err
			}
		}

		filePath := templatesDumpDir + "/" + s.Obj.GetName() + strings.ReplaceAll(s.TemplatePath, "/", "-")
		templateFile, err := os.Create(filePath)
		if err != nil {
			zap.S().Errorf("%v can't create file for rendered template, %v", err, s.Obj.GetName())
			return err
		}
		if _, err = templateFile.Write([]byte(s.ParsedTemplate)); err != nil {
			zap.S().Errorf("%v can't create file for rendered template, %v", err, s.Obj.GetName())
			return err
		}
		if err := templateFile.Close(); err != nil {
			zap.S().Errorf("%v can't close file", err)
			return err
		}

	}
	return nil
}

func GetPromCredsSecret(secretName string, secretNs string, client client.Client, log logr.Logger) (url, user, pass string, err error) {
	user = "cnvrg"
	namespacedName := types.NamespacedName{Name: secretName, Namespace: secretNs}
	creds := v1core.Secret{ObjectMeta: v1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := client.Get(context.Background(), namespacedName, &creds); err != nil && errors.IsNotFound(err) {
		log.Error(err, "Prometheus creds secret not found", "name", secretName)
		return "", "", "", err
	} else if err != nil {
		log.Error(err, "can't get prometheus creds secret", "name", secretName)
		return "", "", "", err
	}

	if _, ok := creds.Data["CNVRG_PROMETHEUS_USER"]; !ok {
		err := fmt.Errorf("prometheus creds secret %s missing require field CNVRG_PROMETHEUS_USER", namespacedName.Name)
		log.Error(err, "missing required field")
		return "", "", "", err
	}

	if _, ok := creds.Data["CNVRG_PROMETHEUS_PASS"]; !ok {
		err := fmt.Errorf("prometheus creds secret %s missing require field CNVRG_PROMETHEUS_PASS", namespacedName.Name)
		log.Error(err, "missing required field")
		return "", "", "", err
	}

	if _, ok := creds.Data["CNVRG_PROMETHEUS_URL"]; !ok {
		err := fmt.Errorf("prometheus creds secret %s missing require field CNVRG_PROMETHEUS_URL", namespacedName.Name)
		log.Error(err, "missing required field")
		return "", "", "", err
	}

	return string(creds.Data["CNVRG_PROMETHEUS_URL"]), string(creds.Data["CNVRG_PROMETHEUS_USER"]), string(creds.Data["CNVRG_PROMETHEUS_PASS"]), nil
}

func GetRedisCredsSecret(secretName string, secretNs string, client client.Client, log logr.Logger) (pass string, err error) {
	namespacedName := types.NamespacedName{Name: secretName, Namespace: secretNs}
	creds := v1core.Secret{ObjectMeta: v1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := client.Get(context.Background(), namespacedName, &creds); err != nil && errors.IsNotFound(err) {
		log.Error(err, "redis creds secret not found", "name", secretName)
		return "", err
	} else if err != nil {
		log.Error(err, "can't get prometheus creds secret", "name", secretName)
		return "", err
	}

	if _, ok := creds.Data["CNVRG_REDIS_PASSWORD"]; !ok {
		err := fmt.Errorf("redis creds secret %s missing require field CNVRG_REDIS_PASSWORD", namespacedName.Name)
		log.Error(err, "missing required field")
		return "", err
	}

	return string(creds.Data["CNVRG_REDIS_PASSWORD"]), nil
}

func CreateRedisCredsSecret(obj v1.Object, secretName, secretNs, redisUrl string, client client.Client, schema *runtime.Scheme, log logr.Logger) error {
	namespacedName := types.NamespacedName{Name: secretName, Namespace: secretNs}
	creds := v1core.Secret{ObjectMeta: v1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := client.Get(context.Background(), namespacedName, &creds); err != nil && errors.IsNotFound(err) {
		if err := ctrl.SetControllerReference(obj, &creds, schema); err != nil {
			log.Error(err, "error set controller reference", "name", namespacedName.Name)
			return err
		}

		pass := RandomString()
		creds.Data = map[string][]byte{
			"CNVRG_REDIS_PASSWORD":              []byte(pass),
			"REDIS_URL":                         []byte(fmt.Sprintf("redis://:%s@%s", pass, redisUrl)), // for cnvrg webapp/sidekiq
			"OAUTH2_PROXY_REDIS_CONNECTION_URL": []byte(fmt.Sprintf("redis://:%s@%s", pass, redisUrl)), // for oauth2 proxy
			"redis.conf":                        []byte(redisConf(pass)),
		}
		if err := client.Create(context.Background(), &creds); err != nil {
			log.Error(err, "error creating redis creds", "name", namespacedName.Name)
			return err
		}
		return nil
	} else if err != nil {
		log.Error(err, "can't check if redis creds secret exists", "name", namespacedName.Name)
		return err
	}
	return nil
}

func CreatePromCredsSecret(obj v1.Object, secretName, secretNs, promUrl string, client client.Client, schema *runtime.Scheme, log logr.Logger) error {
	user := "cnvrg"
	namespacedName := types.NamespacedName{Name: secretName, Namespace: secretNs}
	creds := v1core.Secret{ObjectMeta: v1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := client.Get(context.Background(), namespacedName, &creds); err != nil && errors.IsNotFound(err) {
		if err := ctrl.SetControllerReference(obj, &creds, schema); err != nil {
			log.Error(err, "error set controller reference", "name", namespacedName.Name)
			return err
		}

		pass := RandomString()
		passHash, err := apr1_crypt.New().Generate([]byte(pass), nil)
		if err != nil {
			log.Error(err, "error generating prometheus hash")
			return err
		}
		creds.Data = map[string][]byte{
			"CNVRG_PROMETHEUS_USER": []byte(user),
			"CNVRG_PROMETHEUS_PASS": []byte(pass),
			"CNVRG_PROMETHEUS_URL":  []byte(promUrl),
			"htpasswd":              []byte(fmt.Sprintf("%s:%s", user, passHash)),
		}
		if err := client.Create(context.Background(), &creds); err != nil {
			log.Error(err, "error creating prometheus creds", "name", namespacedName.Name)
			return err
		}

		return nil
	} else if err != nil {
		log.Error(err, "can't check if prometheus creds secret exists", "name", namespacedName.Name)
		return err
	}
	return nil

}

func PrometheusUpstreamConfig(user, pass, ns, upstream string) string {
	return fmt.Sprintf(`
- job_name: 'federate'
  scrape_interval: 10s
  honor_labels: true
  honor_timestamps: false
  metrics_path: '/federate'
  basic_auth:
    username: '%s'
    password: '%s'
  params:
    'match[]':
      - '{namespace="%s"}'
  static_configs:
    - targets:
      - '%s'
`, user, pass, ns, upstream)
}

func redisConf(password string) string {
	return fmt.Sprintf(`
dir /data/
appendonly "yes"
appendfilename "appendonly.aof"
appendfsync everysec
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 128mb
requirepass %s
`, password)
}

func RandomString() string {
	var output strings.Builder
	mathrand.Seed(time.Now().Unix())
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < 20; i++ {
		random := mathrand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

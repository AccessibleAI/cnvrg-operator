package desired

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Masterminds/sprig"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
	"text/template"
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
	return ""
}

func getGrafanaDashboards(obj interface{}) []string {
	if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(obj) {
		return GrafanaInfraDashboards
	}
	if reflect.TypeOf(&mlopsv1.CnvrgApp{}) == reflect.TypeOf(obj) {
		return GrafanaAppDashboards
	}
	return nil
}

func cnvrgTemplateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"ns": func(obj interface{}) string {
			return getNs(obj)
		},
		"httpScheme": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.Networking.HTTPS.Enabled == "true" {
				return "https://"
			}
			return "http://"
		},
		"appDomain": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.NodePortIngress {
				return cnvrgApp.Spec.ClusterDomain + ":" +
					strconv.Itoa(cnvrgApp.Spec.ControlPlan.WebApp.NodePort)
			} else {
				return cnvrgApp.Spec.ControlPlan.WebApp.SvcName + "." + cnvrgApp.Spec.ClusterDomain
			}
		},
		"defaultComputeClusterDomain": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.Networking.Ingress.IngressType == mlopsv1.NodePortIngress {
				return cnvrgApp.Spec.ClusterDomain + ":" +
					strconv.Itoa(cnvrgApp.Spec.ControlPlan.WebApp.NodePort)
			} else {
				return cnvrgApp.Spec.ClusterDomain
			}
		},
		"redisUrl": func(cnvrgApp mlopsv1.CnvrgApp) string {
			return "redis://" + cnvrgApp.Spec.ControlPlan.Redis.SvcName
		},
		"esUrl": func(cnvrgApp mlopsv1.CnvrgApp) string {
			return "http://" + cnvrgApp.Spec.Logging.Es.SvcName
		},
		"hyperServerUrl": func(cnvrgApp mlopsv1.CnvrgApp) string {
			return "http://" + cnvrgApp.Spec.ControlPlan.Hyper.SvcName
		},
		"esFullInternalUrl": func(cnvrgApp mlopsv1.CnvrgApp) string {
			return fmt.Sprintf("http://%s.%s.svc.cluster.local:%d",
				cnvrgApp.Spec.Logging.Es.SvcName,
				cnvrgApp.Namespace,
				cnvrgApp.Spec.Logging.Es.Port)
		},
		"objectStorageUrl": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.ControlPlan.ObjectStorage.CnvrgStorageEndpoint != "" {
				return cnvrgApp.Spec.ControlPlan.ObjectStorage.CnvrgStorageEndpoint
			}
			if cnvrgApp.Spec.Networking.HTTPS.Enabled == "true" {
				return fmt.Sprintf("https://%s.%s", cnvrgApp.Spec.ControlPlan.Minio.SvcName, cnvrgApp.Spec.ClusterDomain)
			} else {
				return fmt.Sprintf("http://%s.%s", cnvrgApp.Spec.ControlPlan.Minio.SvcName, cnvrgApp.Spec.ClusterDomain)
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
		"oauthProxyConfig": func(cnvrgApp mlopsv1.CnvrgApp) string {
			skipAuthUrls := "["
			for i, url := range cnvrgApp.Spec.ControlPlan.OauthProxy.SkipAuthRegex {
				if i == (len(cnvrgApp.Spec.ControlPlan.OauthProxy.SkipAuthRegex) - 1) {
					skipAuthUrls += fmt.Sprintf(`"%v"`, url)
				} else {
					skipAuthUrls += fmt.Sprintf(`"%v", `, url)
				}
			}
			skipAuthUrls += "]"
			proxyConf := []string{
				fmt.Sprintf("provider = %v", cnvrgApp.Spec.ControlPlan.OauthProxy.Provider),
				fmt.Sprintf("http_address = 0.0.0.0:%v", cnvrgApp.Spec.ControlPlan.WebApp.Port),
				fmt.Sprintf("redirect_url = %v", cnvrgApp.Spec.ControlPlan.OauthProxy.RedirectURI),
				fmt.Sprintf("redis_connection_url = redis://%v:%v", cnvrgApp.Spec.ControlPlan.Redis.SvcName, cnvrgApp.Spec.ControlPlan.Redis.Port),
				fmt.Sprintf("redirect_url = %v", cnvrgApp.Spec.ControlPlan.OauthProxy.RedirectURI),
				fmt.Sprintf("skip_auth_regex = %v", skipAuthUrls),
				fmt.Sprintf(`email_domains = ["%v"]`, cnvrgApp.Spec.ControlPlan.OauthProxy.EmailDomain),
				fmt.Sprintf("client_id = %v", cnvrgApp.Spec.ControlPlan.OauthProxy.ClientID),
				fmt.Sprintf("client_secret = %v", cnvrgApp.Spec.ControlPlan.OauthProxy.ClientSecret),
				fmt.Sprintf("cookie_secret = %v", cnvrgApp.Spec.ControlPlan.OauthProxy.CookieSecret),
				fmt.Sprintf("oidc_issuer_url = %v", cnvrgApp.Spec.ControlPlan.OauthProxy.OidcIssuerURL),
				`upstreams = ["http://127.0.0.1:3000/"]`,
				"session_store_type = redis",
				"custom_templates_dir = /opt/app-root/src/templates",
				"ssl_insecure_skip_verify = true",
				"cookie_name = _oauth2_proxy",
				"cookie_expire = 168h",
				"cookie_secure = false",
				"cookie_httponly = true",
			}

			return strings.Join(proxyConf, "\n")
		},
		"cnvrgPassengerBindAddress": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.ControlPlan.OauthProxy.Enabled == "true" {
				return "127.0.0.1"
			}
			return "0.0.0.0"
		},
		"cnvrgPassengerBindPort": func(cnvrgApp mlopsv1.CnvrgApp) int {
			if cnvrgApp.Spec.ControlPlan.OauthProxy.Enabled == "true" {
				return 3000
			}
			return cnvrgApp.Spec.ControlPlan.WebApp.Port
		},
		"prometheusStaticConfig": func(cnvrgApp mlopsv1.CnvrgApp, ns string) string {
			return fmt.Sprintf(`
- job_name: 'federate'
  scrape_interval: 10s
  honor_labels: true
  honor_timestamps: false
  metrics_path: '/federate'
  params:
    'match[]':
      - '{namespace="%s"}'
  static_configs:
    - targets:
      - '%s'
`, ns, cnvrgApp.Spec.Monitoring.UpstreamPrometheus)
		},
		"grafanaDataSource": func(promSvc string, ns string, promPort int) string {
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
            "url": "http://%s.%s.svc:%d",
            "version": 1
        }
    ]
}`, promSvc, ns, promPort)
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
	}
}

func (s *State) GenerateDeployable(spec v1.Object) error {
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
	if err := s.Template.Execute(&tpl, spec); err != nil {
		zap.S().Error(err, "rendering template error", "file", s.TemplatePath)
		return err
	}
	s.ParsedTemplate = tpl.String()
	zap.S().Infof("parsing: %v ", s.TemplatePath)
	zap.S().Debug("template: " + s.TemplatePath + "\n" + s.ParsedTemplate)
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	if _, _, err := dec.Decode([]byte(s.ParsedTemplate), nil, s.Obj); err != nil {
		zap.S().Errorf("%v, template: %v", err, s.ParsedTemplate)
		return err
	}
	return nil
}

func Apply(desiredManifests []*State, desiredSpec v1.Object, client client.Client, schema *runtime.Scheme, log logr.Logger) error {

	ctx := context.Background()
	for _, manifest := range desiredManifests {
		if err := manifest.GenerateDeployable(desiredSpec); err != nil {
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
			log.Info("creating", "name", manifest.Obj.GetName(), "kind", manifest.GVR.Kind)
			if err := client.Create(ctx, manifest.Obj); err != nil {
				log.Error(err, "error creating object", "name", manifest.Obj.GetName())
				return err
			}
		} else {

			if manifest.GVR == Kinds[PvcGVR] {
				// TODO: make this generic
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

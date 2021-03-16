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

func cnvrgTemplateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"httpScheme": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			if cnvrgappspec.Ingress.HTTPS.Enabled == "true" {
				return "https://"
			}
			return "http://"
		},
		"appDomain": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			if cnvrgappspec.Ingress.IngressType == mlopsv1.NodePortIngress {
				return cnvrgappspec.ClusterDomain + ":" +
					strconv.Itoa(cnvrgappspec.ControlPlan.WebApp.NodePort)
			} else {
				return cnvrgappspec.ControlPlan.WebApp.SvcName + "." + cnvrgappspec.ClusterDomain
			}
		},
		"defaultComputeClusterDomain": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			if cnvrgappspec.Ingress.IngressType == mlopsv1.NodePortIngress {
				return cnvrgappspec.ClusterDomain + ":" +
					strconv.Itoa(cnvrgappspec.ControlPlan.WebApp.NodePort)
			} else {
				return cnvrgappspec.ClusterDomain
			}
		},
		"redisUrl": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			return "redis://" + cnvrgappspec.Redis.SvcName
		},
		"esUrl": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			return "http://" + cnvrgappspec.Logging.Es.SvcName
		},
		"hyperServerUrl": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			return "http://" + cnvrgappspec.ControlPlan.Hyper.SvcName
		},
		"objectStorageUrl": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			if cnvrgappspec.ControlPlan.ObjectStorage.CnvrgStorageEndpoint != "" {
				return cnvrgappspec.ControlPlan.ObjectStorage.CnvrgStorageEndpoint
			}
			if cnvrgappspec.Ingress.HTTPS.Enabled == "true" {
				return fmt.Sprintf("https://%s.%s", cnvrgappspec.Minio.SvcName, cnvrgappspec.ClusterDomain)
			} else {
				return fmt.Sprintf("http://%s.%s", cnvrgappspec.Minio.SvcName, cnvrgappspec.ClusterDomain)
			}
		},
		"routeBy": func(cnvrgappspec mlopsv1.CnvrgAppSpec, routeBy string) string {
			switch routeBy {
			case "ISTIO":
				if cnvrgappspec.Ingress.IngressType == mlopsv1.IstioIngress {
					return "true"
				}
				return "false"
			case "OPENSHIFT":
				if cnvrgappspec.Ingress.IngressType == mlopsv1.OpenShiftIngress {
					return "true"
				}
				return "false"
			case "NGINX_INGRESS":
				if cnvrgappspec.Ingress.IngressType == mlopsv1.NginxIngress {
					return "true"
				}
				return "false"
			case "NODE_PORT":
				if cnvrgappspec.Ingress.IngressType == mlopsv1.NodePortIngress {
					return "true"
				}
				return "false"
			}
			return "false"
		},
		"oauthProxyConfig": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			skipAuthUrls := "["
			for i, url := range cnvrgappspec.ControlPlan.OauthProxy.SkipAuthRegex {
				if i == (len(cnvrgappspec.ControlPlan.OauthProxy.SkipAuthRegex) - 1) {
					skipAuthUrls += fmt.Sprintf(`"%v"`, url)
				} else {
					skipAuthUrls += fmt.Sprintf(`"%v", `, url)
				}
			}
			skipAuthUrls += "]"
			proxyConf := []string{
				fmt.Sprintf("provider = %v", cnvrgappspec.ControlPlan.OauthProxy.Provider),
				fmt.Sprintf("http_address = 0.0.0.0:%v", cnvrgappspec.ControlPlan.WebApp.Port),
				fmt.Sprintf("redirect_url = %v", cnvrgappspec.ControlPlan.OauthProxy.RedirectURI),
				fmt.Sprintf("redis_connection_url = redis://%v:%v", cnvrgappspec.Redis.SvcName, cnvrgappspec.Redis.Port),
				fmt.Sprintf("redirect_url = %v", cnvrgappspec.ControlPlan.OauthProxy.RedirectURI),
				fmt.Sprintf("skip_auth_regex = %v", skipAuthUrls),
				fmt.Sprintf(`email_domains = ["%v"]`, cnvrgappspec.ControlPlan.OauthProxy.EmailDomain),
				fmt.Sprintf("client_id = %v", cnvrgappspec.ControlPlan.OauthProxy.ClientID),
				fmt.Sprintf("client_secret = %v", cnvrgappspec.ControlPlan.OauthProxy.ClientSecret),
				fmt.Sprintf("cookie_secret = %v", cnvrgappspec.ControlPlan.OauthProxy.CookieSecret),
				fmt.Sprintf("oidc_issuer_url = %v", cnvrgappspec.ControlPlan.OauthProxy.OidcIssuerURL),
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
		"cnvrgPassengerBindAddress": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			if cnvrgappspec.ControlPlan.OauthProxy.Enabled == "true" {
				return "127.0.0.1"
			}
			return "0.0.0.0"
		},
		"cnvrgPassengerBindPort": func(cnvrgappspec mlopsv1.CnvrgAppSpec) int {
			if cnvrgappspec.ControlPlan.OauthProxy.Enabled == "true" {
				return 3000
			}
			return cnvrgappspec.ControlPlan.WebApp.Port
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
	data:= CastSpec(spec)
	if err := s.Template.Execute(&tpl, &data); err != nil {
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
	s.Name = s.Obj.Object["metadata"].(map[string]interface{})["name"].(string)

	return nil
}

func CastSpec(spec v1.Object) interface{} {
	if reflect.TypeOf(&mlopsv1.CnvrgApp{}) == reflect.TypeOf(spec) {
		cnvrgApp := spec.(*mlopsv1.CnvrgApp)
		return cnvrgApp
	}
	if reflect.TypeOf(&mlopsv1.CnvrgInfra{}) == reflect.TypeOf(spec) {
		cnvrgInfra := spec.(*mlopsv1.CnvrgInfra)
		return cnvrgInfra
	}
	return nil
}

func Apply(desiredManifests []*State, desiredSpec v1.Object, client client.Client, schema *runtime.Scheme, log logr.Logger) error {

	ctx := context.Background()
	for _, manifest := range desiredManifests {
		if err := manifest.GenerateDeployable(desiredSpec); err != nil {
			log.Error(err, "error generating deployable", "name", manifest.Name)
			return err
		}
		if manifest.Own {
			if err := ctrl.SetControllerReference(desiredSpec, manifest.Obj, schema); err != nil {
				log.Error(err, "error setting controller reference", "name", manifest.Name)
				return err
			}
		}
		if viper.GetBool("dry-run") {
			log.Info("dry run enabled, skipping applying...")
			continue
		}
		fetchInto := &unstructured.Unstructured{}
		fetchInto.SetGroupVersionKind(manifest.GVR)
		err := client.Get(ctx, types.NamespacedName{Name: manifest.Name, Namespace: desiredSpec.GetNamespace()}, fetchInto)
		if err != nil && errors.IsNotFound(err) {
			log.Info("creating", "name", manifest.Name, "kind", manifest.GVR.Kind)
			if err := client.Create(ctx, manifest.Obj); err != nil {
				log.Error(err, "error creating object", "name", manifest.Name)
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
			//manifest.Obj.SetResourceVersion(fetchInto.GetResourceVersion())
			err := client.Update(ctx, fetchInto)
			if err != nil {
				log.Info("error updating object", "manifest", manifest.TemplatePath)
				return err
			}
		}
	}
	return nil
}

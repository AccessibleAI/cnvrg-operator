package desired

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/markbates/pkger"
	"go.uber.org/zap"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"strconv"
	"strings"
	"text/template"
)

func cnvrgTemplateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"httpScheme": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			if cnvrgappspec.Networking.HTTPS.Enabled == "true" {
				return "https://"
			}
			return "http://"
		},
		"appDomain": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			if cnvrgappspec.Networking.IngressType == mlopsv1.NodePortIngress {
				return cnvrgappspec.ClusterDomain + ":" +
					strconv.Itoa(cnvrgappspec.ControlPlan.WebApp.NodePort)
			} else {
				return cnvrgappspec.ControlPlan.WebApp.SvcName + "." + cnvrgappspec.ClusterDomain
			}
		},
		"defaultComputeClusterDomain": func(cnvrgappspec mlopsv1.CnvrgAppSpec) string {
			if cnvrgappspec.Networking.IngressType == mlopsv1.NodePortIngress {
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
			if cnvrgappspec.Networking.HTTPS.Enabled == "true" {
				return fmt.Sprintf("https://%s.%s", cnvrgappspec.Minio.SvcName, cnvrgappspec.ClusterDomain)
			} else {
				return fmt.Sprintf("http://%s.%s", cnvrgappspec.Minio.SvcName, cnvrgappspec.ClusterDomain)
			}
		},
		"routeBy": func(cnvrgappspec mlopsv1.CnvrgAppSpec, routeBy string) string {
			switch routeBy {
			case "ISTIO":
				if cnvrgappspec.Networking.IngressType == mlopsv1.IstioIngress {
					return "true"
				}
				return "false"
			case "OPENSHIFT":
				if cnvrgappspec.Networking.IngressType == mlopsv1.OpenShiftIngress {
					return "true"
				}
				return "false"
			case "NGINX_INGRESS":
				if cnvrgappspec.Networking.IngressType == mlopsv1.NginxIngress {
					return "true"
				}
				return "false"
			case "NODE_PORT":
				if cnvrgappspec.Networking.IngressType == mlopsv1.NodePortIngress {
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

func (s *State) GenerateDeployable(cnvrgApp *mlopsv1.CnvrgApp) error {
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
	if err := s.Template.Execute(&tpl, cnvrgApp.Spec); err != nil {
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

package desired

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/Dimss/crypt/apr1_crypt"
	"github.com/Masterminds/sprig"
	yamlgh "github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/golang-jwt/jwt"
	"github.com/imdario/mergo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	appv1 "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var GrafanaAppDashboards = []string{
	"grafana-apiserver.json",
}

var GrafanaInfraDashboards = append([]string{
	"grafana-apiserver.json",
	"grafana-cluster-total.json",
	"grafana-controller-manager.json",
	"grafana-k8s-resources-cluster.json",
	"grafana-k8s-resources-node.json",
	"grafana-kubelet.json",
	"grafana-node-cluster-rsrc-use.json",
	"grafana-node-rsrc-use.json",
	"grafana-nodes.json",
	"grafana-prometheus-remote-write.json",
	"grafana-prometheus.json",
	"grafana-proxy.json",
	"grafana-scheduler.json",
	"grafana-node-exporter.json",
	"grafana-fluentbit.json",
	"grafana-dcgm-exporter.json",
}, GrafanaAppDashboards...)

func getNs(obj interface{}) string {
	if reflect.TypeOf(&mlopsv1.CnvrgThirdParty{}) == reflect.TypeOf(obj) {
		cnvrgInfra := obj.(*mlopsv1.CnvrgThirdParty)
		return cnvrgInfra.Namespace
	}
	if reflect.TypeOf(&mlopsv1.CnvrgApp{}) == reflect.TypeOf(obj) {
		cnvrgApp := obj.(*mlopsv1.CnvrgApp)
		return cnvrgApp.Namespace
	}
	return "cnvrg-infra"
}

func getGrafanaDashboards(obj interface{}) []string {
	if reflect.TypeOf(&mlopsv1.CnvrgThirdParty{}) == reflect.TypeOf(obj) {
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
			if cnvrgApp.Spec.Networking.HTTPS.Enabled {
				return "https://"
			}
			return "http://"
		},
		"appDomain": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.NodePortIngress {
				return cnvrgApp.Spec.ClusterDomain + ":" +
					strconv.Itoa(cnvrgApp.Spec.ControlPlane.WebApp.NodePort)
			} else {
				return cnvrgApp.Spec.ControlPlane.WebApp.SvcName + "." + cnvrgApp.Spec.ClusterDomain
			}
		},
		"defaultComputeClusterDomain": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.NodePortIngress {
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
			return fmt.Sprintf("http://%s.%s.svc.%s:%d",
				cnvrgApp.Spec.Dbs.Es.SvcName,
				cnvrgApp.Namespace,
				cnvrgApp.Spec.ClusterInternalDomain,
				cnvrgApp.Spec.Dbs.Es.Port)
		},
		"objectStorageUrl": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.ControlPlane.ObjectStorage.Endpoint != "" {
				return cnvrgApp.Spec.ControlPlane.ObjectStorage.Endpoint
			} else if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.NodePortIngress {
				return fmt.Sprintf("http://%s:%d", cnvrgApp.Spec.ClusterDomain, cnvrgApp.Spec.Dbs.Minio.NodePort)
			} else {
				if cnvrgApp.Spec.Networking.HTTPS.Enabled {
					return fmt.Sprintf("https://%s.%s", cnvrgApp.Spec.Dbs.Minio.SvcName, cnvrgApp.Spec.ClusterDomain)
				} else {
					return fmt.Sprintf("http://%s.%s", cnvrgApp.Spec.Dbs.Minio.SvcName, cnvrgApp.Spec.ClusterDomain)
				}
			}
		},
		"cnvrgRoutingService": func(cnvrgApp mlopsv1.CnvrgApp) string {
			if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.NodePortIngress {
				return fmt.Sprintf("http://%s:%d", cnvrgApp.Spec.ClusterDomain, cnvrgApp.Spec.ControlPlane.CnvrgRouter.NodePort)
			} else if cnvrgApp.Spec.Networking.HTTPS.Enabled {
				return fmt.Sprintf("https://%s.%s", cnvrgApp.Spec.ControlPlane.CnvrgRouter.SvcName, cnvrgApp.Spec.ClusterDomain)
			} else {
				return fmt.Sprintf("http://%s.%s", cnvrgApp.Spec.ControlPlane.CnvrgRouter.SvcName, cnvrgApp.Spec.ClusterDomain)
			}

		},
		"routeBy": func(cnvrgApp mlopsv1.CnvrgApp, routeBy string) string {
			switch routeBy {
			case "ISTIO":
				if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.IstioIngress {
					return "true"
				}
				return "false"
			case "OPENSHIFT":
				if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.OpenShiftIngress {
					return "true"
				}
				return "false"
			case "NGINX_INGRESS":
				if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.NginxIngress {
					return "true"
				}
				return "false"
			case "NODE_PORT":
				if cnvrgApp.Spec.Networking.Ingress.Type == mlopsv1.NodePortIngress {
					return "true"
				}
				return "false"
			}
			return "false"
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
            "secureJsonData": {
			  "basicAuthPassword": "%s"
            },
            "jsonData": {
                "tlsSkipVerify": true
            }
        }
    ]
}`, promUrl, user, pass)
		},
		"grafanaDashboards": func(obj interface{}) []string {
			return getGrafanaDashboards(obj)
		},
		"isAppSpec": func(obj interface{}) bool {
			if reflect.TypeOf(&mlopsv1.CnvrgThirdParty{}) == reflect.TypeOf(obj) {
				return false
			}
			return true
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
		"isTrue": func(boolPointer bool) bool { // this is legacy function and should be removed in the future
			return boolPointer
		},
		"promRetentionSize": func(retentionSize string) string {
			size, err := strconv.Atoi(strings.TrimSuffix(retentionSize, "Gi"))
			if err != nil {
				zap.S().Error(err)
			}
			return fmt.Sprintf("%dGB", size-2)
		},
		"image": func(imageHub string, imageName string) string {
			if strings.Contains(imageName, "/") {
				return imageName
			} else {
				return fmt.Sprintf("%s/%s", imageHub, imageName)
			}
		},
		"redisConf": func(password string) string {
			passConf := ""
			if password != "" {
				passConf = fmt.Sprintf("requirepass %s", password)
			}
			return fmt.Sprintf(`
dir /data/
appendonly "yes"
appendfilename "appendonly.aof"
appendfsync everysec
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 128mb
timeout 15
%s
`, passConf)
		},
		// token visibility levels: https://github.com/AccessibleAI/metagpu-device-plugin/blob/main/pkg/mgsrv/server.go#L30
		"generateMetagpuToken": func(secret string, tokenLevel string) string {
			claims := jwt.MapClaims{"email": "metagpu@instance", "visibilityLevel": tokenLevel}
			containerScopeToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := containerScopeToken.SignedString([]byte(secret))
			if err != nil {
				fmt.Println(err)
			}
			return tokenString
		},
	}
}

func (s *State) GenerateDeployable() error {
	var tpl bytes.Buffer
	b, err := s.Fs.ReadFile(s.TemplatePath)
	if err != nil {
		zap.S().Error(err, "error reading path", "path", s.TemplatePath)
		return err
	}

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
	s.Obj.SetGroupVersionKind(s.GVK)
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
	return nil
}

func (s *State) mergeMetadata(actualObject *unstructured.Unstructured, log logr.Logger) {
	var e error
	var jsonStr []byte
	if s.GVK.Kind == "Deployment" {

		manifest, err := objectToDeployment(s.Obj)
		if err != nil {
			log.Error(err, "error merging metadata")
			return
		}
		actual, err := objectToDeployment(actualObject)
		if err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		if err := mergo.Merge(&manifest.ObjectMeta.Labels, actual.ObjectMeta.Labels); err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		if err := mergo.Merge(&manifest.Spec.Template.ObjectMeta.Labels, actual.Spec.Template.ObjectMeta.Labels); err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		if err := mergo.Merge(&manifest.ObjectMeta.Annotations, actual.ObjectMeta.Annotations); err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		if err := mergo.Merge(&manifest.Spec.Template.ObjectMeta.Annotations, actual.Spec.Template.ObjectMeta.Annotations); err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		jsonStr, e = json.Marshal(manifest)
	}

	if s.GVK.Kind == "StatefulSet" {

		manifest, err := objectToStatefulSet(s.Obj)
		if err != nil {
			log.Error(err, "error merging metadata")
			return
		}
		actual, err := objectToStatefulSet(actualObject)
		if err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		if err := mergo.Merge(&manifest.ObjectMeta.Labels, actual.ObjectMeta.Labels); err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		if err := mergo.Merge(&manifest.Spec.Template.ObjectMeta.Labels, actual.Spec.Template.ObjectMeta.Labels); err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		if err := mergo.Merge(&manifest.ObjectMeta.Annotations, actual.ObjectMeta.Annotations); err != nil {
			log.Error(err, "error merging metadata")
			return

		}

		if err := mergo.Merge(&manifest.Spec.Template.ObjectMeta.Annotations, actual.Spec.Template.ObjectMeta.Annotations); err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		jsonStr, e = json.Marshal(manifest)

	}

	if s.GVK.Kind == "DaemonSet" {

		manifest, err := objectToDaemonSet(s.Obj)
		if err != nil {
			log.Error(err, "error merging metadata")
			return
		}
		actual, err := objectToDaemonSet(actualObject)
		if err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		if err := mergo.Merge(&manifest.ObjectMeta.Labels, actual.ObjectMeta.Labels); err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		if err := mergo.Merge(&manifest.Spec.Template.ObjectMeta.Labels, actual.Spec.Template.ObjectMeta.Labels); err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		if err := mergo.Merge(&manifest.ObjectMeta.Annotations, actual.ObjectMeta.Annotations); err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		if err := mergo.Merge(&manifest.Spec.Template.ObjectMeta.Annotations, actual.Spec.Template.ObjectMeta.Annotations); err != nil {
			log.Error(err, "error merging metadata")
			return
		}

		jsonStr, e = json.Marshal(manifest)

	}

	if e == nil && len(jsonStr) > 0 {
		dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
		if _, _, err := dec.Decode(jsonStr, nil, s.Obj); err != nil {
			log.Error(err, "error merging metadata")
			return
		}
	} else if e != nil {
		log.Error(e, "error merging metadata")
	}

}

func (s *State) DumpTemplateToFile(preserveTmplDirs bool) error {
	templatesDumpDir := viper.GetString("templates-dump-dir")
	if templatesDumpDir != "" {
		if _, err := os.Stat(templatesDumpDir); os.IsNotExist(err) {
			if err = os.Mkdir(templatesDumpDir, 0775); err != nil {
				zap.S().Error(err, "can't create templates dump dir for templates debugging")
				return err
			}
		}

		filePath := ""
		if preserveTmplDirs {
			filePath = templatesDumpDir + "/" + s.TemplatePath
			if err := os.MkdirAll(filepath.Dir(filePath), 0775); err != nil {
				return err
			}
		} else {
			filePath = templatesDumpDir + "/" + s.Obj.GetName() + strings.ReplaceAll(s.TemplatePath, "/", "-")
		}

		templateFile, err := os.Create(strings.ReplaceAll(filePath, "tpl", "yaml"))
		if err != nil {
			zap.S().Errorf("%v can't create file for rendered template, %v", err, s.Obj.GetName())
			return err
		}
		b, err := s.Obj.MarshalJSON()
		if err != nil {
			return err
		}

		res, err := yamlgh.JSONToYAML(b)
		if err != nil {
			return err
		}

		if _, err = templateFile.Write(res); err != nil {
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

func objectToDeployment(obj interface{}) (*appv1.Deployment, error) {

	d := &appv1.Deployment{}
	if err := objectToUnstructured(obj, d); err != nil {
		return nil, err
	}
	return d, nil
}

func objectToStatefulSet(obj interface{}) (*appv1.StatefulSet, error) {

	s := &appv1.StatefulSet{}
	if err := objectToUnstructured(obj, s); err != nil {
		return nil, err
	}
	return s, nil
}

func objectToDaemonSet(obj interface{}) (*appv1.DaemonSet, error) {

	d := &appv1.DaemonSet{}
	if err := objectToUnstructured(obj, d); err != nil {
		return nil, err
	}
	return d, nil
}

func objectToUnstructured(obj interface{}, dst interface{}) error {
	un, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return err
	}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(un, dst); err != nil {
		return err
	}
	return nil
}

func GetPromCredsSecret(secretName string, secretNs string, client client.Client, log logr.Logger) (url, user, pass string, err error) {
	user = "cnvrg"
	namespacedName := types.NamespacedName{Name: secretName, Namespace: secretNs}
	creds := v1core.Secret{ObjectMeta: v1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}}
	if err := client.Get(context.Background(), namespacedName, &creds); err != nil && errors.IsNotFound(err) {
		log.Error(err, "Prometheus creds secret not found (either not created yet or you are using external prometheus: https://install.cnvrg.io/deployments/openshift.html)", "name", secretName)
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

	return string(creds.Data["CNVRG_PROMETHEUS_URL"]), string(creds.Data["CNVRG_PROMETHEUS_USER"]), string(creds.Data["CNVRG_PROMETHEUS_HASHED_PASS"]), nil
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

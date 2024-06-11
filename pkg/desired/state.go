package desired

import (
	"fmt"
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	mathrand "math/rand"
	"strconv"
	"strings"
	"time"
)

func cnvrgTemplateFuncs() map[string]interface{} {
	return map[string]interface{}{

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
				return cnvrgApp.Spec.ControlPlane.WebApp.SvcName +
					cnvrgApp.Spec.Networking.ClusterDomainPrefix.Prefix + "." +
					cnvrgApp.Spec.ClusterDomain
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
					return fmt.Sprintf("https://%s%s.%s",
						cnvrgApp.Spec.Dbs.Minio.SvcName,
						cnvrgApp.Spec.Networking.ClusterDomainPrefix.Prefix,
						cnvrgApp.Spec.ClusterDomain)
				} else {
					return fmt.Sprintf("http://%s%s.%s",
						cnvrgApp.Spec.Dbs.Minio.SvcName,
						cnvrgApp.Spec.Networking.ClusterDomainPrefix.Prefix,
						cnvrgApp.Spec.ClusterDomain)
				}
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
	}
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

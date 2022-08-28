apiVersion: v1
kind: ConfigMap
metadata:
  name: cp-networking-config
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  APP_DOMAIN: "{{ appDomain . }}"
  DEFAULT_URL: "{{ httpScheme . }}{{ appDomain . }}"
  DEFAULT_COMPUTE_CLUSTER_DOMAIN: "{{ defaultComputeClusterDomain .}}"
  DEFAULT_COMPUTE_CLUSTER_HTTPS: "{{ .Spec.Networking.HTTPS.Enabled }}"
  CNVRG_CLUSTER_INTERNAL_DOMAIN: "{{ .Spec.ClusterInternalDomain }}"
  HYPER_SERVER_URL: "{{ hyperServerUrl .}}"
  HYPER_SERVER_PORT: "{{ .Spec.ControlPlane.Hyper.Port }}"
  ROUTE_BY_ISTIO: "{{ routeBy . "ISTIO" }}"
  ROUTE_BY_OPENSHIFT: "{{ routeBy . "OPENSHIFT" }}"
  ROUTE_BY_NGINX_INGRESS: "{{ routeBy . "NGINX_INGRESS" }}"
  ROUTE_BY_NODE_PORT: "{{ routeBy . "NODE_PORT" }}"
  CNVRG_ISTIO_GATEWAY: "{{ .Spec.Networking.Ingress.IstioGwName }}"
  {{- if isTrue .Spec.ControlPlane.CnvrgRouter.Enabled }}
  DEPLOY_URL: "{{ cnvrgRoutingService . }}"
  NOTEBOOK_URL: "{{ cnvrgRoutingService . }}"
  TENSORBOARD_URL: "{{ cnvrgRoutingService . }}"
  {{- end }}
  {{- if isTrue .Spec.Networking.Proxy.Enabled }}
  CNVRG_PROXY_CONFIG_REF: {{ .Spec.Networking.Proxy.ConfigRef }}
  {{- end }}


apiVersion: v1
kind: ConfigMap
metadata:
  name: cp-networking-config
  namespace: {{ ns . }}
data:
  APP_DOMAIN: "{{ appDomain . }}"
  DEFAULT_URL: "{{ httpScheme . }}{{ appDomain . }}"
  DEFAULT_COMPUTE_CLUSTER_DOMAIN: "{{ defaultComputeClusterDomain .}}"
  DEFAULT_COMPUTE_CLUSTER_HTTPS: "{{ .Spec.Networking.HTTPS.Enabled }}"
  HYPER_SERVER_URL: "{{ hyperServerUrl .}}"
  HYPER_SERVER_PORT: "{{ .Spec.ControlPlane.Hyper.Port }}"
  ROUTE_BY_ISTIO: "{{ routeBy . "ISTIO" }}"
  ROUTE_BY_OPENSHIFT: "{{ routeBy . "OPENSHIFT" }}"
  ROUTE_BY_NGINX_INGRESS: "{{ routeBy . "NGINX_INGRESS" }}"
  ROUTE_BY_NODE_PORT: "{{ routeBy . "NODE_PORT" }}"
  CNVRG_ISTIO_GATEWAY: "{{ istioGwName .}}"



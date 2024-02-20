apiVersion: v1
kind: ConfigMap
metadata:
  name: cp-networking-config
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
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
  CNVRG_SSO_CENTRAL_PROXY: "{{.Spec.SSO.Proxy.Address}}"
  CNVRG_SSO_PUBLIC_URL: "{{.Spec.SSO.Central.PublicUrl}}"
  {{- if isTrue .Spec.Networking.Proxy.Enabled }}
  CNVRG_PROXY_CONFIG_REF: {{ .Spec.Networking.Proxy.ConfigRef }}
  {{- end }}
  {{- if isTrue .Spec.Networking.HTTPS.Enabled }}
  CNVRG_CERT_SECRET: {{ .Spec.Networking.HTTPS.CertSecret }}
  {{- end }}
  {{- if isTrue .Spec.Networking.Ingress.OcpSecureRoutes }}
  OCP_SECURE_ROUTES: "true"
  {{- end }}
  INGRESS_CLASS_NAME: {{ .Spec.Networking.Ingress.IngressClassName }}

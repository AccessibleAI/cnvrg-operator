apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Spec.Networking.Proxy.ConfigRef }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "false"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  http_proxy: {{ .Spec.Networking.Proxy.HttpProxy | join "," }}
  HTTP_PROXY: {{ .Spec.Networking.Proxy.HttpProxy | join "," }}
  https_proxy: {{ .Spec.Networking.Proxy.HttpsProxy | join "," }}
  HTTPS_PROXY: {{ .Spec.Networking.Proxy.HttpsProxy | join "," }}
  no_proxy: {{ .Spec.Networking.Proxy.NoProxy | join "," }}
  NO_PROXY: {{ .Spec.Networking.Proxy.NoProxy | join "," }}
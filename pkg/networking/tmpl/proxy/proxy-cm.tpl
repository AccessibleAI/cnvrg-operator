apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Spec.Networking.Proxy.ConfigRef }}
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
  http_proxy: {{ .Spec.Networking.Proxy.HttpProxy | join "," }}
  HTTP_PROXY: {{ .Spec.Networking.Proxy.HttpProxy | join "," }}
  https_proxy: {{ .Spec.Networking.Proxy.HttpsProxy | join "," }}
  HTTPS_PROXY: {{ .Spec.Networking.Proxy.HttpsProxy | join "," }}
  no_proxy: {{ .Spec.Networking.Proxy.NoProxy | join "," }}
  NO_PROXY: {{ .Spec.Networking.Proxy.NoProxy | join "," }}
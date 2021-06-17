apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Spec.Proxy.ConfigRef }}
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
  http_proxy: {{ .Spec.Proxy.HttpProxy | join "," }}
  HTTP_PROXY: {{ .Spec.Proxy.HttpProxy | join "," }}
  https_proxy: {{ .Spec.Proxy.HttpsProxy | join "," }}
  HTTPS_PROXY: {{ .Spec.Proxy.HttpsProxy | join "," }}
  no_proxy: {{ .Spec.Proxy.NoProxy | join "," }}
  NO_PROXY: {{ .Spec.Proxy.NoProxy | join "," }}
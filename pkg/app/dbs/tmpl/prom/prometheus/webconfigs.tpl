apiVersion: v1
kind: ConfigMap
metadata:
  name: prom-web-configs
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "false"
    {{- range $k, $v := .Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  web-config.yml: |
    basic_auth_users:
      cnvrg: {{ .PassHash }}
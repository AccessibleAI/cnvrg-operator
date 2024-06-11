apiVersion: v1
kind: Secret
metadata:
  name: grafana-datasources
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
type: Opaque
data:
  datasources.yaml: {{ grafanaDataSource .Url .User .Pass | b64enc }}


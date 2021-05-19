apiVersion: v1
kind: Secret
metadata:
  name: grafana-datasources
  namespace: {{ .Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
type: Opaque
data:
  datasources.yaml: {{ grafanaDataSource .Data.Url .Data.User .Data.Pass | b64enc }}


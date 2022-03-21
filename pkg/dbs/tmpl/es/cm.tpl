apiVersion: v1
kind: ConfigMap
metadata:
  name: es-ilm-cm
  namespace: {{ .Data.Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  cleanup_policy_all: '3d'
  cleanup_policy_app: '30d'
  cleanup_policy_jobs: '14d'
  cleanup_policy_endpoints: '1825d'
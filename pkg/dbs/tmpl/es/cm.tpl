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
  CLEANUP_POLICY_ALL: '3d'
  CLEANUP_POLICY_APP: '30d'
  CLEANUP_POLICY_JOBS: '14d'
  CLEANUP_POLICY_ENDPOINTS: '1825d'
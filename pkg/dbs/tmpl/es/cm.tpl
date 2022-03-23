apiVersion: v1
kind: ConfigMap
metadata:
  name: es-ilm-cm
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  CLEANUP_POLICY_ALL: '3d'
  CLEANUP_POLICY_APP: '30d'
  CLEANUP_POLICY_JOBS: '14d'
  CLEANUP_POLICY_ENDPOINTS: '1825d'
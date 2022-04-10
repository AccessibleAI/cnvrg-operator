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
  CLEANUP_POLICY_ALL: {{ .Spec.Dbs.Es.CleanupPolicy.All }}
  CLEANUP_POLICY_APP: {{ .Spec.Dbs.Es.CleanupPolicy.App }}
  CLEANUP_POLICY_JOBS: {{ .Spec.Dbs.Es.CleanupPolicy.Jobs }}
  CLEANUP_POLICY_ENDPOINTS: {{ .Spec.Dbs.Es.CleanupPolicy.Endpoints }}
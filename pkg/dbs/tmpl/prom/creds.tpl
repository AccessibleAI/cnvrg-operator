apiVersion: v1
kind: Secret
metadata:
  name: prom-creds
  namespace: {{ .Data.Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  CNVRG_PROMETHEUS_USER: {{ .Data.User | b64enc }}
  CNVRG_PROMETHEUS_PASS: {{ .Data.Pass | b64enc }}
  CNVRG_PROMETHEUS_HASHED_PASS: {{ .Data.PassHash | b64enc }}
  CNVRG_PROMETHEUS_URL:  {{ .Data.PromUrl | b64enc }}
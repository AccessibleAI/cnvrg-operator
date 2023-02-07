apiVersion: v1
kind: Secret
metadata:
  name: {{ .CredsRef }}
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
  CNVRG_PROMETHEUS_USER: {{ .User | b64enc }}
  CNVRG_PROMETHEUS_PASS: {{ .Pass | b64enc }}
  CNVRG_PROMETHEUS_HASHED_PASS: {{ .PassHash | b64enc }}
  CNVRG_PROMETHEUS_URL: {{ .PromUrl | b64enc }}
  CNVRG_PROMETHEUS_INTERNAL_URL: "{{ .InternalUrl | b64enc }}"
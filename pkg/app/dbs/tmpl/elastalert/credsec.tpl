apiVersion: v1
kind: Secret
metadata:
  name: {{ .Data.CredsRef }}
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
  CNVRG_ELASTALERT_USER: {{ .Data.User | b64enc }}
  CNVRG_ELASTALERT_PASS: {{ .Data.Pass | b64enc }}
  CNVRG_ELASTALERT_URL:  {{ .Data.ElastAlertUrl | b64enc}}
  htpasswd:              {{ .Data.Htpasswd | b64enc }}
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
  CNVRG_ELASTALERT_USER: {{ .User | b64enc }}
  CNVRG_ELASTALERT_PASS: {{ .Pass | b64enc }}
  CNVRG_ELASTALERT_URL:  {{ .ElastAlertUrl | b64enc}}
  htpasswd:              {{ .Htpasswd | b64enc }}
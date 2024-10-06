apiVersion: v1
kind: Secret
metadata:
  name: cp-smtp
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  SMTP_SERVER: {{ .Server }}
  SMTP_PORT: {{ .Port }}
  SMTP_USERNAME: {{ .Username }}
  SMTP_PASSWORD: {{ .Password }}
  SMTP_DOMAIN: {{ .Domain }}
  SMTP_OPENSSL_VERIFY_MODE: {{ .OpensslVerifyMode }}
  SMTP_SENDER: {{ .Sender }}

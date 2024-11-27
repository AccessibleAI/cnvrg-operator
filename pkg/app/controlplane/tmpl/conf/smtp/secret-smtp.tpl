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
    {{- range $k, $v := .Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  SMTP_SERVER: {{ .Server | b64enc }}
  SMTP_PORT: {{ .Port | toString | b64enc }}
  SMTP_USERNAME: {{ .Username | b64enc }}
  SMTP_PASSWORD: {{ .Password | b64enc}}
  SMTP_DOMAIN: {{ .Domain | b64enc}}
  SMTP_OPENSSL_VERIFY_MODE: {{ .OpensslVerifyMode | b64enc }}
  SMTP_SENDER: {{ .Sender | b64enc }}

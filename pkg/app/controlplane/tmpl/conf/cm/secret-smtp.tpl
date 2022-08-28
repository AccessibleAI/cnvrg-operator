apiVersion: v1
kind: Secret
metadata:
  name: cp-smtp
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  SMTP_SERVER: {{ .Spec.ControlPlane.SMTP.Server | b64enc }}
  SMTP_PORT: {{ .Spec.ControlPlane.SMTP.Port | toString | b64enc }}
  SMTP_USERNAME: {{ .Spec.ControlPlane.SMTP.Username | b64enc }}
  SMTP_PASSWORD: {{ .Spec.ControlPlane.SMTP.Password | b64enc }}
  SMTP_DOMAIN: {{ .Spec.ControlPlane.SMTP.Domain | b64enc }}
  SMTP_OPENSSL_VERIFY_MODE: {{ .Spec.ControlPlane.SMTP.OpensslVerifyMode | b64enc }}
  SMTP_SENDER: {{ .Spec.ControlPlane.SMTP.Sender | b64enc }}

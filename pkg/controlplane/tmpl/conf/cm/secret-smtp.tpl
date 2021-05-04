apiVersion: v1
kind: Secret
metadata:
  name: cp-smtp
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
data:
  SMTP_SERVER: {{ .Spec.ControlPlane.SMTP.Server | b64enc }}
  SMTP_PORT: {{ .Spec.ControlPlane.SMTP.Port | toString | b64enc }}
  SMTP_USERNAME: {{ .Spec.ControlPlane.SMTP.Username | b64enc }}
  SMTP_PASSWORD: {{ .Spec.ControlPlane.SMTP.Password | b64enc }}
  SMTP_DOMAIN: {{ .Spec.ControlPlane.SMTP.Domain | b64enc }}

apiVersion: v1
kind: Secret
metadata:
  name: cp-base-secret
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
data:
  SENTRY_URL: {{ .Spec.ControlPlane.BaseConfig.SentryURL | b64enc }}
  HYPER_SERVER_TOKEN: {{ .Spec.ControlPlane.Hyper.Token | b64enc }}
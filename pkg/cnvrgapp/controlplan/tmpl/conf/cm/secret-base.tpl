apiVersion: v1
kind: Secret
metadata:
  name: cp-base-secret
  namespace: {{ .CnvrgNs }}
data:
  SENTRY_URL: {{ .ControlPlan.BaseConfig.SentryURL | b64enc }}
  HYPER_SERVER_TOKEN: {{ .ControlPlan.Hyper.Token | b64enc }}
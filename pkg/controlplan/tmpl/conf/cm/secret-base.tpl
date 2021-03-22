apiVersion: v1
kind: Secret
metadata:
  name: cp-base-secret
  namespace: {{ ns . }}
data:
  SENTRY_URL: {{ .Spec.ControlPlan.BaseConfig.SentryURL | b64enc }}
  HYPER_SERVER_TOKEN: {{ .Spec.ControlPlan.Hyper.Token | b64enc }}
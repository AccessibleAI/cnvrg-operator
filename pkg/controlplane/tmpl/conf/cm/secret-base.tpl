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
  SECRET_KEY_BASE: {{ randAlphaNum 128 | lower | b64enc }}
  STS_IV: {{ "DeJ/CGz/Hkb/IbRe4t1xLg==" | b64enc }}
{{/*  STS_KEY: {{ randAlphaNum 32 | lower | b64enc }}*/}}
  STS_KEY: {{ randAlphaNum 32 | lower | b64enc }}
apiVersion: v1
kind: ConfigMap
metadata:
  name: env-config
  namespace: {{ .CnvrgNs }}
data:
  APP_DOMAIN: {{ appDomain . }}
  DEFAULT_URL: {{ httpScheme . }}{{ appDomain . }}

apiVersion: v1
kind: Secret
metadata:
  name: cp-sso-{{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
data:
  conf: {{ oauthProxyConfig . .Spec.Logging.Kibana.SvcName | b64enc }}


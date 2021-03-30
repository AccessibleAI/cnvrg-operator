apiVersion: v1
kind: Secret
metadata:
  name: "oauth-proxy-{{.Spec.Logging.Kibana.SvcName}}"
  namespace: {{ ns . }}
data:
  conf: {{ oauthProxyConfig . .Spec.Logging.Kibana.SvcName nil | b64enc }}


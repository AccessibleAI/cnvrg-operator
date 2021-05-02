apiVersion: v1
kind: Secret
metadata:
  name: "oauth-proxy-{{.Spec.Logging.Kibana.SvcName}}"
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
data:
  conf: {{ oauthProxyConfig . .Spec.Logging.Kibana.SvcName nil .Spec.SSO.Provider .Spec.Logging.Kibana.Port 3000 | b64enc }}


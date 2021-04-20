apiVersion: v1
kind: Secret
metadata:
  name: "oauth-proxy-{{.Spec.Monitoring.Prometheus.SvcName}}"
  namespace: {{ ns . }}
data:
  conf: {{ oauthProxyConfig . .Spec.Monitoring.Prometheus.SvcName (splitList "&" `^\/static`) "cnvrg" .Spec.Monitoring.Prometheus.Port 9090 | b64enc }}


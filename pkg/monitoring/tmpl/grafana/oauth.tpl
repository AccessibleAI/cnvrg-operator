apiVersion: v1
kind: Secret
metadata:
  name: oauth-proxy-grafana
  namespace: {{ ns . }}
data:
  conf: {{ oauthProxyConfig . .Spec.Monitoring.Grafana.SvcName .Spec.Monitoring.Grafana.OauthProxy.SkipAuthRegex | b64enc }}
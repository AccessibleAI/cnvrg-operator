apiVersion: v1
kind: Secret
metadata:
  name: oauth-proxy-grafana
  namespace: {{ ns . }}
data:
  conf: {{ oauthProxyConfig . .Spec.Monitoring.Grafana.SvcName .Spec.Monitoring.Grafana.OauthProxy.SkipAuthRegex .Spec.SSO.Provider .Spec.Monitoring.Grafana.Port 3000 | b64enc }}
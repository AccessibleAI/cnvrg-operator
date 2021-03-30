apiVersion: v1
kind: Secret
metadata:
  name: oauth-proxy-webapp
  namespace: {{ ns . }}
data:
  conf: {{ oauthProxyConfig . .Spec.ControlPlan.WebApp.SvcName .Spec.ControlPlan.WebApp.OauthProxy.SkipAuthRegex | b64enc }}
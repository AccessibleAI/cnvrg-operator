apiVersion: v1
kind: Secret
metadata:
  name: oauth-proxy-webapp
  namespace: {{ ns . }}
data:
  conf: {{ oauthProxyConfig . .Spec.ControlPlane.WebApp.SvcName .Spec.ControlPlane.WebApp.OauthProxy.SkipAuthRegex | b64enc }}
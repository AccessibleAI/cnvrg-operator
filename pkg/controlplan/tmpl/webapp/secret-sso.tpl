apiVersion: v1
kind: Secret
metadata:
  name: cp-sso-app
  namespace: {{ ns . }}
data:
  conf: {{ oauthProxyConfig . .Spec.ControlPlan.WebApp.SvcName .Spec.ControlPlan.WebApp.OauthProxy.SkipAuthRegex | b64enc }}
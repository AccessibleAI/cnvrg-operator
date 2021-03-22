apiVersion: v1
kind: Secret
metadata:
  name: cp-sso
  namespace: {{ .Namespace }}
data:
  conf: {{ oauthProxyConfig . | b64enc }}
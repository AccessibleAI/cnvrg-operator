apiVersion: v1
kind: Secret
metadata:
  name: cp-sso
  namespace: {{ ns . }}
data:
  conf: {{ oauthProxyConfig . | b64enc }}
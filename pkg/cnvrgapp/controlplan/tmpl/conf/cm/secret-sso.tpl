apiVersion: v1
kind: Secret
metadata:
  name: cp-sso
  namespace: {{ .CnvrgNs }}
data:
  conf: {{ oauthProxyConfig . | b64enc }}
apiVersion: v1
kind: Secret
metadata:
  name: cp-oauth-proxy-tokens-secret
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "false"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  OAUTH_PROXY_API_KEY: {{ randAlphaNum 32 | lower | b64enc }}
  OAUTH_PROXY_API_AUTH_DATA: {{ randAlphaNum 12 | lower | b64enc }}

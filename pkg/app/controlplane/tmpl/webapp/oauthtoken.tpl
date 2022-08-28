apiVersion: v1
kind: Secret
metadata:
  name: cp-oauth-proxy-tokens-secret
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  OAUTH_PROXY_API_KEY: {{ randAlphaNum 32 | lower | b64enc }}
  OAUTH_PROXY_API_AUTH_DATA: {{ randAlphaNum 12 | lower | b64enc }}

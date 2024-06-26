apiVersion: v1
kind: Secret
metadata:
  name: cp-base-secret
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
  CNVRG_SSO_KEY: {{ printf "%s:%s" .Spec.SSO.Central.ClientID .Spec.SSO.Central.ClientSecret | b64enc }}
  OAUTH_PROXY_TOKENS_ENABLED: "{{ .Spec.SSO.Enabled | toString | b64enc }}"
  SENTRY_URL: {{ .Spec.ControlPlane.BaseConfig.SentryURL | b64enc }}
  HYPER_SERVER_TOKEN: {{ .Spec.ControlPlane.Hyper.Token | b64enc }}
  SECRET_KEY_BASE: {{ randAlphaNum 128 | lower | b64enc }}
  STS_IV: {{ "DeJ/CGz/Hkb/IbRe4t1xLg==" | b64enc }}
{{/*  STS_KEY: {{ randAlphaNum 32 | lower | b64enc }} -> waiting for big boss */}}
  STS_KEY: {{ "05646d3cbf8baa5be7150b4283eda07d" | b64enc }}

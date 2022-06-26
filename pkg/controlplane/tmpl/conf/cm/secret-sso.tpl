apiVersion: v1
kind: Secret
metadata:
  name: cp-sso
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
  CNVRG_SSO_KEY: {{ printf "%s:%s" .Spec.SSO.ClientID .Spec.SSO.ClientSecret | b64enc }}
  SSO_IDP_PRIVATE_KEY_REF: {{ .Spec.Pki.PrivateKeySecret }}
  OAUTH_PROXY_TOKENS_ENABLED: "{{ .Spec.SSO.Enabled | toString | b64enc }}"
  OAUTH_PROXY_ENABLED: "{{ isTrue .Spec.SSO.Enabled }}"
  OAUTH_ADMIN_USER: "{{ .Spec.SSO.AdminUser }}"
  CNVRG_SSO_REALM: "{{ .Spec.SSO.RealmName }}"
  CNVRG_SSO_SERVICE_URL: "{{ .Spec.SSO.ServiceUrl }}"
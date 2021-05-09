{{- define "spec.sso" }}
sso:
  enabled: {{ .Values.sso.enabled}}
  adminUser: "{{ .Values.sso.adminUser }}"
  provider: "{{ .Values.sso.provider }}"
  emailDomain:
  {{- range $_, $value := .Values.sso.emailDomain }}
    - {{$value}}
  {{- end }}
  clientId: "{{ .Values.sso.clientId }}"
  clientSecret: "{{ .Values.sso.clientSecret }}"
  azureTenant: "{{ .Values.sso.azureTenant }}"
  oidcIssuerUrl: "{{ .Values.sso.oidcIssuerUrl }}"
  cookieSecret: {{ randAlphaNum 16 | quote }}
{{- end }}
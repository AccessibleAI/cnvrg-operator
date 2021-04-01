{{- define "spec.sso" }}
sso:
  enabled: "{{ .Values.sso.enabled}}"
  image: {{ .Values.sso.image }}
  redisConnectionUrl: {{ .Values.sso.redisConnectionUrl }}
  adminUser: "{{ .Values.sso.adminUser }}"
  provider: "{{ .Values.sso.provider }}"
  emailDomain: "{{ .Values.sso.emailDomain }}"
  clientId: "{{ .Values.sso.clientId }}"
  clientSecret: "{{ .Values.sso.clientSecret }}"
  cookieSecret: "{{ .Values.sso.cookieSecret }}"
  azureTenant: "{{ .Values.sso.azureTenant }}"
  oidcIssuerUrl: "{{ .Values.sso.oidcIssuerURL }}"
{{- end }}
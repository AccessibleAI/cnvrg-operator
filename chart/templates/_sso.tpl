sso:
  enabled: {{ .Values.sso.enabled}}
  image: {{ .Values.sso.image }}
  redisConnectionUrl: {{ .Values.sso.redisConnectionUrl }}
  adminUser: {{ .Values.sso.adminUser }}
  provider: {{ .Values.sso.provider }}
  emailDomain: {{ .Values.sso.emailDomain }}
  clientID: {{ .Values.sso.clientID }}
  clientSecret: {{ .Values.sso.clientSecret }}
  cookieSecret: {{ .Values.sso.cookieSecret }}
  azureTenant: {{ .Values.sso.azureTenant }}
  oidcIssuerURL: {{ .Values.sso.oidcIssuerURL }}
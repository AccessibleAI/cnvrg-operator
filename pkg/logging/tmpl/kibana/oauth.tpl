apiVersion: v1
kind: Secret
metadata:
  name: "oauth-proxy-{{.Spec.Logging.Kibana.SvcName}}"
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-kibana-oauth"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  conf: {{ oauthProxyConfig . .Spec.Logging.Kibana.SvcName nil .Spec.SSO.Provider .Spec.Logging.Kibana.Port 3000 .Spec.ControlPlane.WebApp.OauthProxy.TokenValidationKey .Spec.ControlPlane.WebApp.OauthProxy.TokenValidationAuthData .Spec.ControlPlane.WebApp.OauthProxy.TokenValidationRegex | b64enc }}


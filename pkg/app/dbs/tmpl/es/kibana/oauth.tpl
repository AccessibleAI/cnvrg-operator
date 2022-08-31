apiVersion: v1
kind: Secret
metadata:
  name: "oauth-proxy-{{.Spec.Dbs.Es.Kibana.SvcName}}"
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-kibana-oauth"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  conf: {{ oauthProxyConfig . .Spec.Dbs.Es.Kibana.SvcName nil .Spec.SSO.Provider .Spec.Dbs.Es.Kibana.Port 3000 .Spec.Dbs.Es.Kibana.OauthProxy.TokenValidationRegex | b64enc }}


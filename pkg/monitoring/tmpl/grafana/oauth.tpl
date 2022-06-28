apiVersion: v1
kind: Secret
metadata:
  name: oauth-proxy-grafana
  namespace: {{ ns . }}
  annotations:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-grafana-oauth"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  conf: {{ oauthProxyConfig . .Spec.Monitoring.Grafana.SvcName .Spec.Monitoring.Grafana.OauthProxy.SkipAuthRegex .Spec.SSO.Provider .Spec.Monitoring.Grafana.Port 3000 .Spec.Monitoring.Grafana.OauthProxy.TokenValidationRegex | b64enc }}
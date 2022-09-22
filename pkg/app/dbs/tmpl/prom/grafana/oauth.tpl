apiVersion: v1
kind: Secret
metadata:
  name: oauth-proxy-grafana
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
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
  conf: {{ oauthProxyConfig . .Spec.Dbs.Prom.Grafana.SvcName .Spec.Dbs.Prom.Grafana.OauthProxy.SkipAuthRegex .Spec.SSO.Provider .Spec.Dbs.Prom.Grafana.Port 3000 .Spec.Dbs.Prom.Grafana.OauthProxy.TokenValidationRegex | b64enc }}
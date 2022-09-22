apiVersion: v1
kind: Secret
metadata:
  name: cp-sso-grafana
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  conf: {{ oauthProxyConfig . .Spec.Dbs.Prom.Grafana.SvcName | b64enc }}


apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.Dbs.Prom.Grafana.CredsRef }}
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
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
  CNVRG_GRAFANA_SVC_NAME: {{ .Spec.Dbs.Prom.Grafana.SvcName | b64enc }}
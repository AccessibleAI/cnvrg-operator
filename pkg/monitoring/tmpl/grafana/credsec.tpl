apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.Monitoring.Grafana.CredsRef }}
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
  CNVRG_GRAFANA_SVC_NAME: {{ .Spec.Monitoring.Grafana.SvcName | b64enc }}
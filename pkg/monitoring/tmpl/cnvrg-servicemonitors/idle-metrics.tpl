apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cnvrg-jobs
  namespace: {{ ns . }}
  labels:
    app: cnvrg-jobs
    cnvrg-ccp-prometheus: {{ .Name }}-{{ ns .}}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
    {{- range $name, $value := .Spec.Monitoring.CnvrgIdleMetricsExporter.Labels }}
    {{ $name }}: {{ $value }}
    {{- end }}
spec:
  jobLabel: cnvrg-job
  selector:
    matchLabels:
      exporter: cnvrg-job
  namespaceSelector:
    matchNames:
      - {{ ns . }}
  endpoints:
    - interval: 30s
      scrapeTimeout: 10s
      port: "http"
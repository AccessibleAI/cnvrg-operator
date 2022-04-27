apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cnvrg-exporter
  namespace: {{ ns . }}
  labels:
    app: cnvrg-exporter
    cnvrg-ccp-prometheus: {{ .Name }}-{{ ns .}}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
    {{- range $name, $value := .Spec.Monitoring.CnvrgUsageMetricsExporter.Labels }}
    {{ $name }}: {{ $value }}
    {{- end }}
spec:
  jobLabel: cnvrg-exporter
  selector:
    matchLabels:
      exporter: cnvrg-exporter
  namespaceSelector:
    matchNames:
      - {{ ns . }}
  endpoints:
    - interval: 30s
      scrapeTimeout: 10s
      port: "http"

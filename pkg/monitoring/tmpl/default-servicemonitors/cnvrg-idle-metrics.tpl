apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cnvrg-jobs
  namespace: {{ ns . }}
  labels:
    app: cnvrg-jobs
    {{- if isTrue .Spec.NamespaceTenancy }}
    cnvrg-ccp-prometheus: {{ .Name }}-{{ ns .}}
    {{- else }}
    cnvrg-infra-prometheus: {{ .Name }}-{{ ns .}}
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
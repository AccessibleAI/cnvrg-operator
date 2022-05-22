apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: habana-exporter
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app.kubernetes.io/name: habana-exporter
    app.kubernetes.io/version: v0.0.1
    app: "habana-exporter"
    cnvrg-infra-prometheus: {{ .Name }}-{{ ns .}}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: habana-exporter
  endpoints:
  - port: habana-metrics
    relabelings:
    - action: replace
      sourceLabels:
      - exported_namespace
      targetLabel: namespace
    interval: 30s
    scrapeTimeout: 20s
  namespaceSelector:
    matchNames:
      - {{ ns . }}

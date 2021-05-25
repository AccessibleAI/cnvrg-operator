apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cnvrg-fluentbit-exporter
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: fluentbit
    cnvrg-infra-prometheus: {{ .Name }}-{{ ns .}}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      app: cnvrg-fluentbit
  namespaceSelector:
    matchNames:
      - {{ ns . }}
  endpoints:
    - port: "metrics"
      path: "/api/v1/metrics/prometheus"
      interval: "15s"
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: metagpu-exporter
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: "metagpu-exporter"
    cnvrg-infra-prometheus: {{ .Name }}-{{ ns .}}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    matchLabels:
      app: "metagpu-exporter"
  namespaceSelector:
    matchNames:
      - {{ ns . }}
  endpoints:
    - port: "metrics"
      path: "/metrics"
      interval: "15s"
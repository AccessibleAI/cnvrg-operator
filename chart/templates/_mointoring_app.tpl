{{- define "spec.monitoring_app" }}
{{- if eq .Values.spec "ccp" }}
monitoring:
  grafana:
    enabled: {{ .Values.monitoring.grafana.enabled }}
  prometheus:
    enabled: {{ .Values.monitoring.prometheus.enabled }}
    storageClass: "{{ .Values.monitoring.prometheus.storageClass }}"
    storageSize: {{ .Values.monitoring.prometheus.storageSize }}
    nodeSelector:
    {{- range $key, $value := .Values.monitoring.prometheus.nodeSelector }}
      {{$key}}: {{$value}}
    {{- end }}
  cnvrgIdleMetricsExporter:
    enabled: {{ .Values.monitoring.cnvrgIdleMetricsExporter.enabled }}
{{- end }}
{{- end }}
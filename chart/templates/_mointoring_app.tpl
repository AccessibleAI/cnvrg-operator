{{- define "spec.monitoring_app" }}
{{- if eq .Values.spec "ccp" }}
monitoring:
  grafana:
    enabled: {{ .Values.monitoring.grafana.enabled }}
    svcName: {{ .Values.monitoring.grafana.svcName }}
  prometheus:
    enabled: {{ .Values.monitoring.prometheus.enabled }}
    storageClass: "{{ .Values.monitoring.prometheus.storageClass }}"
    retention: {{ .Values.monitoring.prometheus.retention }}
    storageSize: {{ .Values.monitoring.prometheus.storageSize }}
    {{- if .Values.monitoring.prometheus.nodeSelector }}
    nodeSelector:
    {{- range $key, $value := .Values.monitoring.prometheus.nodeSelector }}
      {{$key}}: {{$value}}
    {{- end }}
    {{- end }}
  cnvrgIdleMetricsExporter:
    enabled: {{ .Values.monitoring.cnvrgIdleMetricsExporter.enabled }}
{{- end }}
{{- end }}
{{- define "spec.monitoring_infra" }}
monitoring:
  defaultServiceMonitors:
    enabled: {{ .Values.monitoring.defaultServiceMonitors.enabled }}
  {{- if eq .Values.spec "allinone" }}
  cnvrgIdleMetricsExporter:
    enabled: {{ .Values.monitoring.cnvrgIdleMetricsExporter.enabled }}
    labels:
      cnvrg-infra-prometheus: cnvrg-infra-{{ template "spec.cnvrgNs" . }}
    {{- range $key, $value := .Values.monitoring.cnvrgIdleMetricsExporter.labels }}
      {{$key}}: {{$value}}
    {{- end }}
  {{- end }}
  dcgmExporter:
    enabled: {{ .Values.monitoring.dcgmExporter.enabled }}
  grafana:
    enabled: {{ .Values.monitoring.grafana.enabled }}
  kubeStateMetrics:
    enabled: {{ .Values.monitoring.kubeStateMetrics.enabled }}
  nodeExporter:
    enabled: {{ .Values.monitoring.nodeExporter.enabled }}
  prometheus:
    enabled: {{ .Values.monitoring.prometheus.enabled }}
    storageClass: "{{ .Values.monitoring.prometheus.storageClass }}"
    storageSize: {{ .Values.monitoring.prometheus.storageSize }}
    {{- if .Values.monitoring.prometheus.nodeSelector }}
    nodeSelector:
    {{- range $key, $value := .Values.monitoring.prometheus.nodeSelector }}
      {{$key}}: {{$value}}
    {{- end }}
    {{- end }}
  prometheusOperator:
    enabled: {{ .Values.monitoring.prometheusOperator.enabled }}
{{- end }}
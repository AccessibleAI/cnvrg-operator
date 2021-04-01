{{- define "spec.monitoring_app" }}
monitoring:
  enabled: "{{ .Values.monitoring.enabled }}"
  upstreamPrometheus: {{ .Values.monitoring.upstreamPrometheus }}

  grafana:
    enabled: "{{ .Values.monitoring.grafana.enabled }}"
    image: {{ .Values.monitoring.grafana.image }}
    nodePort: {{ .Values.monitoring.grafana.nodePort }}
    port: {{ .Values.monitoring.grafana.port }}
    svcName: {{ .Values.monitoring.grafana.svcName }}
    oauthProxy:
      skipAuthRegex:
        - \/api\/health

  prometheus:
    cpuRequest: {{ .Values.monitoring.prometheus.cpuRequest }}
    enabled: "{{ .Values.monitoring.prometheus.enabled }}"
    image: {{ .Values.monitoring.prometheus.image }}
    memoryRequest: {{ .Values.monitoring.prometheus.memoryRequest }}
    nodePort: {{ .Values.monitoring.prometheus.nodePort }}
    port: {{ .Values.monitoring.prometheus.port }}
    storageClass: "{{ .Values.monitoring.prometheus.storageClass }}"
    storageSize: {{ .Values.monitoring.prometheus.storageSize }}
    svcName: {{ .Values.monitoring.prometheus.svcName }}

{{- end }}
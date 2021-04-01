{{- define "spec.monitoring_infra" }}
monitoring:
  enabled: "{{ .Values.monitoring.enabled }}"
  kubeletServiceMonitor: "{{ .Values.monitoring.kubeletServiceMonitor }}"
  dcgmExporter:
    enabled: "{{ .Values.monitoring.dcgmExporter.enabled }}"
    image: {{ .Values.monitoring.dcgmExporter.image }}
  grafana:
    enabled: "{{ .Values.monitoring.grafana.enabled }}"
    image: {{ .Values.monitoring.grafana.image }}
    nodePort: {{ .Values.monitoring.grafana.nodePort }}
    port: {{ .Values.monitoring.grafana.port }}
    svcName: {{ .Values.monitoring.grafana.svcName }}
    oauthProxy:
      skipAuthRegex:
        - \/api\/health
  kubeStateMetrics:
    enabled: "{{ .Values.monitoring.kubeStateMetrics.enabled }}"
    image: {{ .Values.monitoring.kubeStateMetrics.image }}
  nodeExporter:
    enabled: "{{ .Values.monitoring.nodeExporter.enabled }}"
    image: {{ .Values.monitoring.nodeExporter.image }}
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
  prometheusOperator:
    enabled: "{{ .Values.monitoring.prometheusOperator.enabled }}"
    images:
      kubeRbacProxyImage: {{ .Values.monitoring.prometheusOperator.images.kubeRbacProxyImage }}
      operatorImage: {{ .Values.monitoring.prometheusOperator.images.operatorImage }}
      prometheusConfigReloaderImage: {{ .Values.monitoring.prometheusOperator.images.prometheusConfigReloaderImage }}
{{- end }}
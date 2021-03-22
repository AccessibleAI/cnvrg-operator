apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: cnvrg-infra-prometheus
  namespace: {{ ns . }}
  labels:
    app: cnvrg-infra-prometheus
spec:
  storage:
    disableMountSubPath: true
    volumeClaimTemplate:
      spec:
        resources:
          requests:
            storage: {{ .Spec.Monitoring.Prometheus.StorageSize }}
        {{- if ne .Spec.Monitoring.Prometheus.StorageClass "use-default" }}
        storageClassName: {{ .Spec.Monitoring.Prometheus.StorageClass }}
        {{- end }}
  image: {{ .Spec.Monitoring.Prometheus.Image }}
  replicas: 1
  resources:
    requests:
      cpu: {{ .Spec.Monitoring.Prometheus.CPURequest }}
      memory: {{ .Spec.Monitoring.Prometheus.MemoryRequest }}
  ruleSelector:
    matchLabels:
      app: cnvrg-infra-prometheus
      role: alert-rules
  securityContext:
    fsGroup: 2000
    runAsNonRoot: true
    runAsUser: 1000
  serviceAccountName: cnvrg-infra-prometheus
  podMonitorNamespaceSelector: {}
  podMonitorSelector: {}
  probeNamespaceSelector: {}
  serviceMonitorNamespaceSelector: {}
  serviceMonitorSelector: {}
  version: v2.22.1

apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: cnvrg-ccp-prometheus
  namespace: {{ .Namespace }}
  labels:
    app: cnvrg-ccp-prometheus
spec:
  storage:
    disableMountSubPath: true
    volumeClaimTemplate:
      spec:
        resources:
          requests:
            storage: {{ .Spec.Prometheus.StorageSize }}
        {{- if ne .Spec.Prometheus.StorageClass "use-default" }}
        storageClassName: {{ .Spec.Prometheus.StorageClass }}
        {{- end }}
  image: {{ .Spec.Prometheus.Image }}
  replicas: 1
  resources:
    requests:
      cpu: {{ .Spec.Prometheus.CPURequest }}
      memory: {{ .Spec.Prometheus.MemoryRequest }}
  ruleSelector:
    matchLabels:
      app: cnvrg-ccp-prometheus
      role: alert-rules
  securityContext:
    fsGroup: 2000
    runAsNonRoot: true
    runAsUser: 1000
  serviceAccountName: cnvrg-ccp-prometheus
  podMonitorNamespaceSelector:
    matchLabels:
      namespace: {{ .Namespace }}
  podMonitorSelector:
    matchLabels:
      namespace: {{ .Namespace }}
  probeNamespaceSelector:
    matchLabels:
      namespace: {{ .Namespace }}
  serviceMonitorNamespaceSelector:
    matchLabels:
      namespace: {{ .Namespace }}
  serviceMonitorSelector:
    matchLabels:
      namespace: {{ .Namespace }}
  version: v2.22.1
  additionalScrapeConfigs:
    name: static-config
    key: prometheus-additional.yaml

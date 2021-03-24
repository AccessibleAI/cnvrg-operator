apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: cnvrg-ccp-prometheus
  namespace: {{ ns . }}
  labels:
    app: cnvrg-ccp-prometheus
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
  replicaExternalLabelName: ""
  prometheusExternalLabelName: ""
  replicas: 1
  resources:
    requests:
      cpu: {{ .Spec.Monitoring.Prometheus.CPURequest }}
      memory: {{ .Spec.Monitoring.Prometheus.MemoryRequest }}
  ruleSelector:
    matchLabels:
      app: cnvrg-ccp-prometheus
      role: alert-rules
  securityContext:
    fsGroup: 2000
    runAsNonRoot: true
    runAsUser: 1000
  serviceAccountName: cnvrg-ccp-prometheus
  serviceMonitorNamespaceSelector: {}
  podMonitorNamespaceSelector: {}
  probeNamespaceSelector: {}
  serviceMonitorSelector:
    matchLabels:
      cnvrg-infra-prometheus: "bla"
  podMonitorSelector: {}
  version: v2.22.1
  additionalScrapeConfigs:
    name: prom-static-config
    key: prometheus-additional.yaml

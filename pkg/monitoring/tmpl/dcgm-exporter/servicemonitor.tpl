apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: dcgm-exporter
  namespace: {{ ns . }}
  labels:
    app: "dcgm-exporter"
    owner: cnvrg-control-plane
spec:
  selector:
    matchLabels:
      app: "dcgm-exporter"
  namespaceSelector:
    matchNames:
      - {{ ns . }}
  endpoints:
    - port: "metrics"
      path: "/metrics"
      interval: "15s"
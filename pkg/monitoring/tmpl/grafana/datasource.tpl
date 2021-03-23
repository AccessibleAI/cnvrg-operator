apiVersion: v1
kind: Secret
metadata:
  name: grafana-datasources
  namespace: {{ ns . }}
type: Opaque
data:
  datasources.yaml: {{ grafanaDataSource .Spec.Monitoring.Prometheus.SvcName (ns .) .Spec.Monitoring.Prometheus.Port  | b64enc }}


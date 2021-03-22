apiVersion: v1
kind: Secret
metadata:
  name: grafana-datasources
  namespace: {{ .Namespace }}
type: Opaque
data:
  datasources.yaml: {{ grafanaDataSource . | b64enc }}


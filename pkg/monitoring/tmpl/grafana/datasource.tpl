apiVersion: v1
kind: Secret
metadata:
  name: grafana-datasources
  namespace: {{ .Namespace }}
type: Opaque
data:
  datasources.yaml: {{ grafanaDataSource .Data.Svc .Namespace .Data.Port .Data.User .Data.Pass | b64enc }}


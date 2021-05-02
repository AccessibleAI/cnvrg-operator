apiVersion: v1
kind: Secret
metadata:
  name: grafana-datasources
  namespace: {{ .Namespace }}
  labels:
    owner: cnvrg-control-plane
type: Opaque
data:
  datasources.yaml: {{ grafanaDataSource .Data.Url .Data.User .Data.Pass | b64enc }}


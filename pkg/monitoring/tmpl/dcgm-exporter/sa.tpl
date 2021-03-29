apiVersion: v1
kind: ServiceAccount
metadata:
  name: dcgm-exporter
  namespace: {{ ns . }}
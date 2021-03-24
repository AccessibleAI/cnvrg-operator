apiVersion: v1
kind: ServiceAccount
metadata:
  name: cnvrg-prometheus-operator
  labels:
    app: cnvrg-prometheus-operator
    version: v0.44.1
  namespace: {{ ns . }}

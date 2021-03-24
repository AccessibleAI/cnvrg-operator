apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: {{ ns . }}
  name: istio-operator
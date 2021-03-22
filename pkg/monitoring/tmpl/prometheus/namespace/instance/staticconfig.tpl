apiVersion: v1
kind: Secret
metadata:
  name: static-config
  namespace: {{ .Namespace }}
data:
  prometheus-additional.yaml: {{ prometheusStaticConfig . | b64enc }}
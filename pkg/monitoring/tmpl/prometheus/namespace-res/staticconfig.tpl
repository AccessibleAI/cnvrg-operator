apiVersion: v1
kind: Secret
metadata:
  name: static-config
  namespace: {{ ns . }}
data:
  prometheus-additional.yaml: {{ prometheusStaticConfig . | b64enc }}
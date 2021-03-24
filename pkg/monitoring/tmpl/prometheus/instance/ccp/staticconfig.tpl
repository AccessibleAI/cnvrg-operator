apiVersion: v1
kind: Secret
metadata:
  name: prom-static-config
  namespace: {{ ns . }}
data:
  prometheus-additional.yaml: {{ prometheusStaticConfig . (ns .) | b64enc }}
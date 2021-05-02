apiVersion: v1
kind: Secret
metadata:
  name: kibana-config
  namespace: {{ .Namespace }}
  labels:
    owner: cnvrg-control-plane
data:
  kibana.yml: {{ kibanaSecret .Data.Host .Data.Port .Data.EsHost .Data.EsUser .Data.EsPass (printf "%s:%s" .Data.EsUser .Data.EsPass | b64enc) | b64enc }}

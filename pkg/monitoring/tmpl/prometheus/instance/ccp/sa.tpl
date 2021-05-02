apiVersion: v1
kind: ServiceAccount
metadata:
  name: cnvrg-ccp-prometheus
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}

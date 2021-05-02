apiVersion: v1
kind: ServiceAccount
metadata:
  name: node-exporter
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: dcgm-exporter
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
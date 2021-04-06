apiVersion: v1
kind: ServiceAccount
metadata:
  name: node-exporter
  namespace: {{ ns . }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
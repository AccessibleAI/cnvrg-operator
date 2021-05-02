apiVersion: v1
kind: ServiceAccount
metadata:
  name: fluent-bit
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
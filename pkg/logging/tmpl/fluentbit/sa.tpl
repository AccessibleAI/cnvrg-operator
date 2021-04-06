apiVersion: v1
kind: ServiceAccount
metadata:
  name: fluent-bit
  namespace: {{ ns . }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
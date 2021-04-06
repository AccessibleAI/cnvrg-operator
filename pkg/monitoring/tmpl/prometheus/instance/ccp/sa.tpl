apiVersion: v1
kind: ServiceAccount
metadata:
  name: cnvrg-ccp-prometheus
  namespace: {{ ns . }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
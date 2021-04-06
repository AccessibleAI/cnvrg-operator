apiVersion: v1
kind: ServiceAccount
metadata:
  name: cnvrg-prometheus-operator
  namespace: {{ ns . }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}

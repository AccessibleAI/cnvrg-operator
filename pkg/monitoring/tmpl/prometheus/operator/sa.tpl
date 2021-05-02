apiVersion: v1
kind: ServiceAccount
metadata:
  name: cnvrg-prometheus-operator
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}

apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: {{ .Namespace }}
  name: istio-operator
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
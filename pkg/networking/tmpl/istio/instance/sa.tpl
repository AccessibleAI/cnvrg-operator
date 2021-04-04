apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: {{ ns . }}
  name: istio-operator
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
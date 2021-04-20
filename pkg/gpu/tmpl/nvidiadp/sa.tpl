apiVersion: v1
kind: ServiceAccount
metadata:
  name: nvidia-device-plugin
  namespace: {{ .Namespace }}
imagePullSecrets:
  - name: {{ .Data.Registry.Name }}
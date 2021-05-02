apiVersion: v1
kind: ServiceAccount
metadata:
  name: nvidia-device-plugin
  namespace: {{ .Namespace }}
  labels:
    owner: cnvrg-control-plane
imagePullSecrets:
  - name: {{ .Data.Registry.Name }}
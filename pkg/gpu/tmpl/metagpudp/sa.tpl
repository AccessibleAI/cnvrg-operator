apiVersion: v1
kind: ServiceAccount
metadata:
  name: metagpu-device-plugin
  namespace: {{ .Namespace }}
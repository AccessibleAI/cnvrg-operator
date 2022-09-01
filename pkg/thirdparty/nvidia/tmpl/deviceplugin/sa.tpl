apiVersion: v1
kind: ServiceAccount
metadata:
  name: nvidia-device-plugin
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "false"
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
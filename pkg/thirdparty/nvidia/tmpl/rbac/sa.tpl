apiVersion: v1
kind: ServiceAccount
metadata:
  name: nvidia
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
apiVersion: v1
kind: ConfigMap
metadata:
  name: metagpu-presence
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
data:
  enabled: "{{ .Spec.ControlPlane.BaseConfig.MetagpuEnabled }}"
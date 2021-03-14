apiVersion: v1
kind: ServiceAccount
metadata:
  name: fluentd
  namespace: {{ .CnvrgNs }}
imagePullSecrets:
  - name: {{ .ControlPlan.Registry.Name }}
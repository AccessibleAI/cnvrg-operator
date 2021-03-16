apiVersion: v1
kind: ServiceAccount
metadata:
  name: hostpath-provisioner-admin
  namespace: {{ .CnvrgNs }}
imagePullSecrets:
  - name: {{ .ControlPlan.Registry.Name }}
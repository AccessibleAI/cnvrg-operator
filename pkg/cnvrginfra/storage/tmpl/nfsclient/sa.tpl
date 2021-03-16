apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-client-provisioner
  namespace: {{ .CnvrgNs }}
imagePullSecrets:
  - name: {{ .ControlPlan.Registry.Name }}
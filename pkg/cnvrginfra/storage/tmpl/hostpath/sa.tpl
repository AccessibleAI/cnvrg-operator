apiVersion: v1
kind: ServiceAccount
metadata:
  name: hostpath-provisioner-admin
  namespace: {{ .Spec.CnvrgInfraNs }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
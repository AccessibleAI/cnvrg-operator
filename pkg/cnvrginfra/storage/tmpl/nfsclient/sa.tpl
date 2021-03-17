apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-client-provisioner
  namespace: {{ .Spec.CnvrgInfraNs }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
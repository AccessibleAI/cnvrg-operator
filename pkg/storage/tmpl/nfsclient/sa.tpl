apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-client-provisioner
  namespace: {{ .Namespace }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
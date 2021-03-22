apiVersion: v1
kind: ServiceAccount
metadata:
  name: hostpath-provisioner-admin
  namespace: {{ .Namespace }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
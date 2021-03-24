apiVersion: v1
kind: ServiceAccount
metadata:
  name: hostpath-provisioner-admin
  namespace: {{ ns . }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
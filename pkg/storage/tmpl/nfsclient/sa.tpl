apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-client-provisioner
  namespace: {{ ns . }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
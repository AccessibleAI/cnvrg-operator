apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: {{ .Storage.Nfs.StorageClassName }}
  namespace: {{ .CnvrgNs }}
  annotations:
    storageclass.kubernetes.io/is-default-class: "{{ .Storage.Nfs.DefaultSc }}"
provisioner: {{ .Storage.Nfs.Provisioner }}
reclaimPolicy: {{ .Storage.Nfs.ReclaimPolicy }}
{{- if eq .Storage.Nfs.ReclaimPolicy "Delete" }}
parameters:
  archiveOnDelete: "false"
{{- end }}
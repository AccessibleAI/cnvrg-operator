apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: {{ .Spec.Storage.Nfs.StorageClassName }}
  annotations:
    storageclass.kubernetes.io/is-default-class: "{{ .Spec.Storage.Nfs.DefaultSc }}"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
provisioner: {{ .Spec.Storage.Nfs.Provisioner }}
reclaimPolicy: {{ .Spec.Storage.Nfs.ReclaimPolicy }}
{{- if eq .Spec.Storage.Nfs.ReclaimPolicy "Delete" }}
parameters:
  archiveOnDelete: "false"
{{- end }}
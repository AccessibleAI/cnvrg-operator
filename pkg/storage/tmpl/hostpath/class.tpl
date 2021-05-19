apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: {{ .Spec.Storage.Hostpath.StorageClassName }}
  annotations:
    storageclass.kubernetes.io/is-default-class: "{{ .Spec.Storage.Hostpath.DefaultSc }}"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
provisioner: kubevirt.io/hostpath-provisioner
volumeBindingMode: WaitForFirstConsumer
reclaimPolicy: {{ .Spec.Storage.Hostpath.ReclaimPolicy }}
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: {{ .Spec.Storage.Hostpath.StorageClassName }}
  annotations:
    storageclass.kubernetes.io/is-default-class: "{{ .Spec.Storage.Hostpath.DefaultSc }}"
provisioner: kubevirt.io/hostpath-provisioner
volumeBindingMode: WaitForFirstConsumer
reclaimPolicy: {{ .Spec.Storage.Hostpath.ReclaimPolicy }}
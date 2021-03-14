apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: {{ .Storage.Hostpath.StorageClassName }}
  annotations:
    storageclass.kubernetes.io/is-default-class: {{ .Storage.Hostpath.DefaultSc }}
provisioner: kubevirt.io/hostpath-provisioner
volumeBindingMode: WaitForFirstConsumer
reclaimPolicy: {{ .Storage.Hostpath.ReclaimPolicy }}
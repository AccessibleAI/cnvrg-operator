apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Pg.SvcName }}
  namespace: {{ .Spec.CnvrgNs }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Spec.Pg.StorageSize }}
        {{- if ne .Spec.Pg.StorageClass "use-default" }}
        storageClassName: {{ .Spec.Pg.StorageClass }}
        {{- else if ne .Spec.Storage.CcpStorageClass "" }}
        storageClassName: {{ .Spec.Storage.CcpStorageClass }}
        {{- end }}

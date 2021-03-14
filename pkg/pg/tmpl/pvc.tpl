apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Pg.SvcName }}
  namespace: {{ .CnvrgNs }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Pg.StorageSize }}
  {{- if ne .Pg.StorageClass "use-default" }}
  storageClassName: {{ .Pg.StorageClass }}
  {{- else if ne .Storage.CcpStorageClass "" }}
  storageClassName: {{ .Storage.CcpStorageClass }}
  {{- end }}

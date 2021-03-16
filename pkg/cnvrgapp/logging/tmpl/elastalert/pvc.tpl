apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Logging.Elastalert.SvcName}}
  namespace: {{ .CnvrgNs }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Logging.Elastalert.StorageSize }}
  {{- if ne .Pg.StorageClass "use-default" }}
  storageClassName: {{ .Logging.Elastalert.StorageClass }}
  {{- else if ne .Storage.CcpStorageClass "" }}
  storageClassName: {{ .Storage.CcpStorageClass }}
  {{- end }}
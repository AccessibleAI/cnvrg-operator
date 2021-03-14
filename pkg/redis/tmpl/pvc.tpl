apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name:  {{.Redis.SvcName}}
  namespace: {{ .CnvrgNs }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Redis.StorageSize }}
        {{- if ne .Redis.StorageClass "use-default" }}
        storageClassName: {{ .Redis.StorageClass }}
        {{- else if ne .Storage.CcpStorageClass "" }}
        storageClassName: {{ .Storage.CcpStorageClass }}
        {{- end }}

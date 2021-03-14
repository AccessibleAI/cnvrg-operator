apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Minio.SvcName }}
  namespace: {{ .CnvrgNs  }}
spec:
  accessModes:
    {{- if eq .Minio.SharedStorage.Enabled "true" }}
    - ReadWriteMany
    {{- else }}
    - ReadWriteOnce
    {{- end }}
  resources:
    requests:
      storage: {{ .Minio.StorageSize}}
      {{- if ne .Minio.StorageClass "use-default" }}
      storageClassName: {{ .Minio.StorageClass }}
      {{- else if ne .Storage.CcpStorageClass "" }}
      storageClassName: {{ .Storage.CcpStorageClass }}
      {{- end }}
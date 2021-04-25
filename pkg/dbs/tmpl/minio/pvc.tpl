apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{ ns .  }}
spec:
  accessModes:
    {{- if .Spec.Dbs.Minio.SharedStorage.Enabled }}
    - ReadWriteMany
    {{- else }}
    - ReadWriteOnce
    {{- end }}
  resources:
    requests:
      storage: {{ .Spec.Dbs.Minio.StorageSize}}
  {{- if ne .Spec.Dbs.Minio.StorageClass "" }}
  storageClassName: {{ .Spec.Dbs.Minio.StorageClass }}
  {{- end }}
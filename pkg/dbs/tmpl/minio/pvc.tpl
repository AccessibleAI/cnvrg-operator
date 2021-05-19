apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Dbs.Minio.PvcName }}
  namespace: {{ ns .  }}
  labels:
    owner: cnvrg-control-plane
spec:
  accessModes:
    {{- if isTrue .Spec.Dbs.Minio.SharedStorage.Enabled }}
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
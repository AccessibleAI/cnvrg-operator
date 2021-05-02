apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name:  {{.Spec.Dbs.Redis.SvcName}}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Spec.Dbs.Redis.StorageSize }}
  {{- if ne .Spec.Dbs.Redis.StorageClass "" }}
  storageClassName: {{ .Spec.Dbs.Redis.StorageClass }}
  {{- end }}

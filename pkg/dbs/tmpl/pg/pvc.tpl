apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Dbs.Pg.SvcName }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Spec.Dbs.Pg.StorageSize }}
  {{- if ne .Spec.Dbs.Pg.StorageClass "" }}
  storageClassName: {{ .Spec.Dbs.Pg.StorageClass }}
  {{- end }}

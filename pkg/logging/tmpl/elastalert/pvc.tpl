apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Logging.Elastalert.SvcName}}
  namespace: {{ ns . }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Spec.Logging.Elastalert.StorageSize }}
  {{- if ne .Spec.Logging.Elastalert.StorageClass "" }}
  storageClassName: {{ .Spec.Logging.Elastalert.StorageClass }}
  {{- end }}
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Logging.Elastalert.PvcName }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Spec.Logging.Elastalert.StorageSize }}
  {{- if ne .Spec.Logging.Elastalert.StorageClass "" }}
  storageClassName: {{ .Spec.Logging.Elastalert.StorageClass }}
  {{- end }}
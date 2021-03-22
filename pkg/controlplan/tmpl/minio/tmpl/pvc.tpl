apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Minio.SvcName }}
  namespace: {{ ns .  }}
spec:
  accessModes:
    {{- if eq .Spec.Minio.SharedStorage.Enabled "true" }}
    - ReadWriteMany
    {{- else }}
    - ReadWriteOnce
    {{- end }}
  resources:
    requests:
      storage: {{ .Spec.Minio.StorageSize}}
  {{- if ne .Spec.Minio.StorageClass "use-default" }}
  storageClassName: {{ .Spec.Minio.StorageClass }}
  {{- else if ne .Spec.ControlPlan.BaseConfig.CcpStorageClass "" }}
  storageClassName: {{ .Spec.ControlPlan.BaseConfig.CcpStorageClass }}
  {{- end }}
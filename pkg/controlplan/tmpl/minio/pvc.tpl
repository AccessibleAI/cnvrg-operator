apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.ControlPlan.Minio.SvcName }}
  namespace: {{ ns .  }}
spec:
  accessModes:
    {{- if eq .Spec.ControlPlan.Minio.SharedStorage.Enabled "true" }}
    - ReadWriteMany
    {{- else }}
    - ReadWriteOnce
    {{- end }}
  resources:
    requests:
      storage: {{ .Spec.ControlPlan.Minio.StorageSize}}
  {{- if ne .Spec.ControlPlan.Minio.StorageClass "use-default" }}
  storageClassName: {{ .Spec.ControlPlan.Minio.StorageClass }}
  {{- else if ne .Spec.ControlPlan.BaseConfig.CcpStorageClass "" }}
  storageClassName: {{ .Spec.ControlPlan.BaseConfig.CcpStorageClass }}
  {{- end }}
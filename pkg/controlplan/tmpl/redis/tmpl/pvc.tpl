apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name:  {{.Spec.Redis.SvcName}}
  namespace: {{ ns . }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Spec.Redis.StorageSize }}
  {{- if ne .Spec.Redis.StorageClass "use-default" }}
  storageClassName: {{ .Spec.Redis.StorageClass }}
  {{- else if ne .Spec.ControlPlan.BaseConfig.CcpStorageClass "" }}
  storageClassName: {{ .Spec.ControlPlan.BaseConfig.CcpStorageClass }}
  {{- end }}

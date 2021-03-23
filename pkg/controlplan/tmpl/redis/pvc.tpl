apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name:  {{.Spec.ControlPlan.Redis.SvcName}}
  namespace: {{ ns . }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Spec.ControlPlan.Redis.StorageSize }}
  {{- if ne .Spec.ControlPlan.Redis.StorageClass "use-default" }}
  storageClassName: {{ .Spec.ControlPlan.Redis.StorageClass }}
  {{- else if ne .Spec.ControlPlan.BaseConfig.CcpStorageClass "" }}
  storageClassName: {{ .Spec.ControlPlan.BaseConfig.CcpStorageClass }}
  {{- end }}

apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.ControlPlan.Pg.SvcName }}
  namespace: {{ ns . }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Spec.ControlPlan.Pg.StorageSize }}
  {{- if ne .Spec.ControlPlan.Pg.StorageClass "use-default" }}
  storageClassName: {{ .Spec.ControlPlan.Pg.StorageClass }}
  {{- else if ne .Spec.ControlPlan.BaseConfig.CcpStorageClass "" }}
  storageClassName: {{ .Spec.Spec.ControlPlan.BaseConfig.CcpStorageClass }}
  {{- end }}

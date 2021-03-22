apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Pg.SvcName }}
  namespace: {{ ns . }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Spec.Pg.StorageSize }}
  {{- if ne .Spec.Pg.StorageClass "use-default" }}
  storageClassName: {{ .Spec.Pg.StorageClass }}
  {{- else if ne .Spec.ControlPlan.BaseConfig.CcpStorageClass "" }}
  storageClassName: {{ .Spec.Spec.ControlPlan.BaseConfig.CcpStorageClass }}
  {{- end }}

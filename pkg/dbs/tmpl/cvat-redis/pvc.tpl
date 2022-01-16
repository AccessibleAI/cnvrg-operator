apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Dbs.Cvat.Redis.PvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Spec.Dbs.Cvat.Redis.StorageSize }}
  {{- if ne .Spec.Dbs.Cvat.Redis.StorageClass "" }}
  storageClassName: {{ .Spec.Dbs.Cvat.Redis.StorageClass }}
  {{- end }}

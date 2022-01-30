apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Dbs.Cvat.Pg.PvcName }}
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
      storage: {{ .Spec.Dbs.Cvat.Pg.StorageSize }}
  {{- if ne .Spec.Dbs.Cvat.Pg.StorageClass "" }}
  storageClassName: {{ .Spec.Dbs.Cvat.Pg.StorageClass }}
  {{- end }}

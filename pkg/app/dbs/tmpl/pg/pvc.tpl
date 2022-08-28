apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Dbs.Pg.PvcName }}
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
      storage: {{ .Spec.Dbs.Pg.StorageSize }}
  {{- if ne .Spec.Dbs.Pg.StorageClass "" }}
  storageClassName: {{ .Spec.Dbs.Pg.StorageClass }}
  {{- end }}

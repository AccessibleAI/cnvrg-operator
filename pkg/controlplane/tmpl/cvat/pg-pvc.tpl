apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.ControlPlane.Cvat.Pg.PvcName }}
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
      storage: {{ .Spec.ControlPlane.Cvat.Pg.StorageSize }}
  {{- if ne .Spec.ControlPlane.Cvat.Pg.StorageClass "" }}
  storageClassName: {{ .Spec.ControlPlane.Cvat.Pg.StorageClass }}
  {{- end }}

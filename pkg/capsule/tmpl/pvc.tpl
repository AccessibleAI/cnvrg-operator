apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Capsule.SvcName }}
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
      storage: {{ .Spec.Capsule.StorageSize }}
  {{- if ne .Spec.Capsule.StorageClass "" }}
  storageClassName: {{ .Spec.Capsule.StorageClass }}
  {{- end }}
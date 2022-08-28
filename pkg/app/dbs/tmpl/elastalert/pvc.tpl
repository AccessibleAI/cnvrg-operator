apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Logging.Elastalert.PvcName }}
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
      storage: {{ .Spec.Logging.Elastalert.StorageSize }}
  {{- if ne .Spec.Logging.Elastalert.StorageClass "" }}
  storageClassName: {{ .Spec.Logging.Elastalert.StorageClass }}
  {{- end }}
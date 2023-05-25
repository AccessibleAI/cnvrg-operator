apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Spec.Dbs.Es.Elastalert.PvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "false"
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
      storage: {{ .Spec.Dbs.Es.Elastalert.StorageSize }}
  {{- if ne .Spec.Dbs.Es.Elastalert.StorageClass "" }}
  storageClassName: {{ .Spec.Dbs.Es.Elastalert.StorageClass }}
  {{- end }}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: prom
  namespace: {{ .Data.Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
imagePullSecrets:
  - name: {{ .Data.RegistryName }}
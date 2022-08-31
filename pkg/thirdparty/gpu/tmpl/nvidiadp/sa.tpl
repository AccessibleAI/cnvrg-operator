apiVersion: v1
kind: ServiceAccount
metadata:
  name: nvidia-device-plugin
  namespace: {{ .Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
imagePullSecrets:
  - name: {{ .Data.Registry.Name }}
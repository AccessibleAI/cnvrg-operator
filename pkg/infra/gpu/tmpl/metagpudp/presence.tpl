apiVersion: v1
kind: ConfigMap
metadata:
  name: metagpu-presence
  namespace: {{ .Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  enabled: "{{ .Data.Enabled }}"
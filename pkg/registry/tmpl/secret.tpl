apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  name: {{ .Data.Registry.Name }}
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
  .dockerconfigjson: {{ printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}" .Data.Registry.URL .Data.Registry.User .Data.Registry.Password  (printf "%s:%s" .Data.Registry.User .Data.Registry.Password | b64enc) | b64enc }}


apiVersion: v1
kind: Secret
metadata:
  name: kibana-config
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
  kibana.yml: {{ kibanaSecret .Data.Host .Data.Port .Data.EsHost .Data.EsUser .Data.EsPass (printf "%s:%s" .Data.EsUser .Data.EsPass | b64enc) | b64enc }}

apiVersion: v1
kind: Secret
metadata:
  name: kibana-config
  namespace: {{ .Namespace }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  kibana.yml: {{ kibanaSecret .Data.Host .Data.Port .Data.EsHost .Data.EsUser .Data.EsPass (printf "%s:%s" .Data.EsUser .Data.EsPass | b64enc) | b64enc }}

apiVersion: v1
kind: Secret
metadata:
  name: kibana-config
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  kibana.yml: {{ kibanaSecret .Host .Port .EsHost .EsUser .EsPass (printf "%s:%s" .EsUser .EsPass | b64enc) | b64enc }}

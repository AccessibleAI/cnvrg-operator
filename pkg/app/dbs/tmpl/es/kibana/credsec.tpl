apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.Logging.Kibana.CredsRef }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  CNVRG_KIBANA_SVC_NAME: {{ .Spec.Logging.Kibana.SvcName | b64enc }}
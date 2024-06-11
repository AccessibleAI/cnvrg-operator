apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.Dbs.Es.Kibana.CredsRef }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  CNVRG_KIBANA_SVC_NAME: {{ .Spec.Dbs.Es.Kibana.SvcName | b64enc }}
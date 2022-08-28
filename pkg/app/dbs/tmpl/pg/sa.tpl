apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Spec.Dbs.Pg.ServiceAccount }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
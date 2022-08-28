apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Dbs.Pg.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: cnvrg-postgres
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  ports:
    - port: {{ .Spec.Dbs.Pg.Port }}
  selector:
    app: {{ .Spec.Dbs.Pg.SvcName }}
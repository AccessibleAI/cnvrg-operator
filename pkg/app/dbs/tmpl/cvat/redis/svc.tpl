apiVersion: v1
kind: Service
metadata:
  name: {{.Spec.Dbs.Cvat.Redis.SvcName}}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{.Spec.Dbs.Cvat.Redis.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  ports:
  - name: redis
    port: {{ .Spec.Dbs.Cvat.Redis.Port }}
  selector:
    app: {{ .Spec.Dbs.Cvat.Redis.SvcName }}
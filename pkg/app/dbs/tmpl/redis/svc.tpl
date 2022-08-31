apiVersion: v1
kind: Service
metadata:
  name: {{.Spec.Dbs.Redis.SvcName}}
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{.Spec.Dbs.Redis.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  ports:
  - name: redis
    port: {{ .Spec.Dbs.Redis.Port }}
  selector:
    app: {{ .Spec.Dbs.Redis.SvcName }}
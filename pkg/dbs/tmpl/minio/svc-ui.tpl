apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Dbs.Minio.SvcName }}-ui
  namespace: {{ ns .  }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Minio.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  ports:
  - name: http
    port: 80
    targetPort: 9090
  selector:
    app: {{ .Spec.Dbs.Minio.SvcName }}

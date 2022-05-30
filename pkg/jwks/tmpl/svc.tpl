apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Jwks.Name }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: cnvrg-jwks
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  ports:
    - port: 8080
  selector:
    app: {{ .Spec.Jwks.Name }}
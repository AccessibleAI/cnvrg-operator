apiVersion: v1
kind: Service
metadata:
  name: cnvrg-jwks
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
    - port: 80
  selector:
    app: cnvrg-jwks
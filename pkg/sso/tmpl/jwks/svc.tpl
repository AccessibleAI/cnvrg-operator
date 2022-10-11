apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.SSO.Jwks.Name }}
  namespace: {{.Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
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
    app: {{ .Spec.SSO.Jwks.Name }}
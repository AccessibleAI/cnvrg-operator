apiVersion: v1
kind: Service
metadata:
  name: cnvrg-prometheus-operator
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: cnvrg-prometheus-operator
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  ports:
  - name: https
    port: 8443
    targetPort: https
  selector:
    app: cnvrg-prometheus-operator

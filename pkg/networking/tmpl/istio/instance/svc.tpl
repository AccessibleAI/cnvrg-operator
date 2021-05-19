apiVersion: v1
kind: Service
metadata:
  namespace:  {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    name: istio-operator
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: istio-operator
spec:
  ports:
    - name: http-metrics
      port: 8383
      targetPort: 8383
  selector:
    name: istio-operator
apiVersion: v1
kind: Service
metadata:
  name: cnvrg-fluentbit-exporter
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: cnvrg-fluentbit
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  type: ClusterIP
  ports:
    - name: metrics
      port: 2020
      targetPort: 2020
      protocol: TCP
  selector:
    app: cnvrg-fluentbit
apiVersion: v1
kind: Service
metadata:
  name: fluentbit-exporter
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: fluentbit
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
    app: fluentbit
apiVersion: v1
kind: Service
metadata:
  name: dcgm-exporter
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: dcgm-exporter
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  type: ClusterIP
  ports:
    - name: "metrics"
      port: 9400
      targetPort: 9400
      protocol: TCP
  selector:
    app: "dcgm-exporter"
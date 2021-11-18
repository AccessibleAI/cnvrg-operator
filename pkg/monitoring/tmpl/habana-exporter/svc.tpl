apiVersion: v1
kind: Service
metadata:
  name: habana-exporter
  namespace: {{ ns . }}
  annotations:
    alpha.monitoring.coreos.com/non-namespaced: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app.kubernetes.io/name: habana-exporter
    app.kubernetes.io/version: v0.0.1
    app: habana-exporter
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  clusterIP: None
  ports:
  - name: habana-metrics
    port: 41611
  selector:
    app.kubernetes.io/name: habana-exporter
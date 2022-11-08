apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Dbs.Prom.Grafana.SvcName }}
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Prom.Grafana.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - name: http
    port: {{ .Spec.Dbs.Prom.Grafana.Port }}
    targetPort: {{ .Spec.Dbs.Prom.Grafana.Port }}
    {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
    nodePort: {{ .Spec.Dbs.Prom.Grafana.NodePort }}
    {{- end }}
  selector:
    app: {{ .Spec.Dbs.Prom.Grafana.SvcName }}

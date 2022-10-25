apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Dbs.Es.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ .Spec.Dbs.Es.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - port: {{ .Spec.Dbs.Es.Port }}
    name: http
    {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
    nodePort: {{ .Spec.Dbs.Es.NodePort }}
    {{- end }}
  - name: transport
    protocol: TCP
    port: 9300
    {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
    nodePort: 9300
    {{- end }}
  selector:
    app: {{ .Spec.Dbs.Es.SvcName }}
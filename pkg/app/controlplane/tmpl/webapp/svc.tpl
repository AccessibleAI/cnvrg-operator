apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
  namespace: {{ ns . }}
  annotations:
   {{- range $k, $v := .Spec.Annotations }}
   {{$k}}: "{{$v}}"
   {{- end }}
  labels:
    app: {{ .Spec.ControlPlane.WebApp.SvcName }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - port: {{.Spec.ControlPlane.WebApp.Port}}
    {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
    nodePort: {{ .Spec.ControlPlane.WebApp.NodePort }}
    {{- end }}
  selector:
    app: {{ .Spec.ControlPlane.WebApp.SvcName }}
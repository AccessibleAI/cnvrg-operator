apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ .Spec.ControlPlane.WebApp.SvcName }}
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
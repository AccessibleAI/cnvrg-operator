apiVersion: v1
kind: Service
metadata:
  name: {{ .ControlPlan.WebApp.SvcName }}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{ .ControlPlan.WebApp.SvcName }}
spec:
  {{- if eq .Networking.IngressType "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - port: {{.ControlPlan.WebApp.Port}}
    {{- if eq .Networking.IngressType "nodeport" }}
    nodePort: {{ .ControlPlan.WebApp.NodePort }}
    {{- end }}
  selector:
    app: {{ .ControlPlan.WebApp.SvcName }}
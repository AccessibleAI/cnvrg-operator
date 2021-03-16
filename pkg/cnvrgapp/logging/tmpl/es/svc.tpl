
apiVersion: v1
kind: Service
metadata:
  name: {{ .Logging.Es.SvcName }}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{ .Logging.Es.SvcName }}
spec:
  {{- if eq .Networking.IngressType "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - port: {{ .Logging.Es.Port}}
    {{- if eq .Networking.IngressType "nodeport" }}
    nodePort: {{ .Logging.Es.NodePort }}
    {{- end }}
  selector:
    app: {{ .Logging.Es.SvcName }}

apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Logging.Es.SvcName }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Spec.Logging.Es.SvcName }}
spec:
  {{- if eq .Spec.Ingress.IngressType "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - port: {{ .Spec.Logging.Es.Port}}
    {{- if eq .Spec.Ingress.IngressType "nodeport" }}
    nodePort: {{ .Spec.Logging.Es.NodePort }}
    {{- end }}
  selector:
    app: {{ .Spec.Logging.Es.SvcName }}
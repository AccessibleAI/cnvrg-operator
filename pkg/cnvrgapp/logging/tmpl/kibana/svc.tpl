
apiVersion: v1
kind: Service
metadata:
  name: {{ .Logging.Kibana.SvcName }}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{ .Logging.Kibana.SvcName }}
spec:
  {{- if eq .Networking.IngressType "nodeport" }}
  type: NodePort
  {{- end }}
  selector:
    app: {{ .Logging.Kibana.SvcName }}
  ports:
    - port: {{ .Logging.Kibana.Port }}
      protocol: TCP
      {{- if eq .Networking.IngressType "nodeport" }}
      nodePort: {{ .Logging.Kibana.NodePort }}
      {{- end }}
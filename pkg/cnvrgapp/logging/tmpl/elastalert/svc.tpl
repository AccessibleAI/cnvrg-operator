
apiVersion: v1
kind: Service
metadata:
  name: {{ .Logging.Elastalert.SvcName }}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{ .Logging.Elastalert.SvcName }}
spec:
  {{- if eq .Networking.IngressType "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
    - port: {{ .Logging.Elastalert.Port }}
      protocol: TCP
      targetPort: {{ .Logging.Elastalert.ContainerPort}}
      {{- if eq .Networking.IngressType "nodeport" }}
      nodePort: {{ .Logging.Elastalert.NodePort }}
      {{- end }}
  selector:
    app: {{ .Logging.Elastalert.SvcName }}
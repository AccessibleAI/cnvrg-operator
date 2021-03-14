apiVersion: v1
kind: Service
metadata:
  name: {{ .Minio.SvcName }}
  namespace: {{ .CnvrgNs  }}
  labels:
    app: {{ .Minio.SvcName }}
spec:
  {{- if eq .Networking.IngressType "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - name: http
    port: 80
    targetPort: {{ .Minio.Port }}
    {{- if eq .Networking.IngressType "nodeport" }}
    nodePort: {{ .Minio.NodePort }}
    {{- end }}
  selector:
    app: {{ .Minio.SvcName }}
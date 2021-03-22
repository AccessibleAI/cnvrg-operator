apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Minio.SvcName }}
  namespace: {{ .Namespace  }}
  labels:
    app: {{ .Spec.Minio.SvcName }}
spec:
  {{- if eq .Spec.Ingress.IngressType "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - name: http
    port: 80
    targetPort: {{ .Spec.Minio.Port }}
    {{- if eq .Spec.Ingress.IngressType "nodeport" }}
    nodePort: {{ .Spec.Minio.NodePort }}
    {{- end }}
  selector:
    app: {{ .Spec.Minio.SvcName }}
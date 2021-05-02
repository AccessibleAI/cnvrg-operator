apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{ ns .  }}
  labels:
    app: {{ .Spec.Dbs.Minio.SvcName }}
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - name: http
    port: 80
    targetPort: {{ .Spec.Dbs.Minio.Port }}
    {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
    nodePort: {{ .Spec.Dbs.Minio.NodePort }}
    {{- end }}
  selector:
    app: {{ .Spec.Dbs.Minio.SvcName }}
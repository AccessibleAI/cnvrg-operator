apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.ControlPlan.Minio.SvcName }}
  namespace: {{ ns .  }}
  labels:
    app: {{ .Spec.ControlPlan.Minio.SvcName }}
spec:
  {{- if eq .Spec.Networking.Ingress.IngressType "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - name: http
    port: 80
    targetPort: {{ .Spec.ControlPlan.Minio.Port }}
    {{- if eq .Spec.Networking.Ingress.IngressType "nodeport" }}
    nodePort: {{ .Spec.ControlPlan.Minio.NodePort }}
    {{- end }}
  selector:
    app: {{ .Spec.ControlPlan.Minio.SvcName }}
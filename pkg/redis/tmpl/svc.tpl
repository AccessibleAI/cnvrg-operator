apiVersion: v1
kind: Service
metadata:
  name: {{.Redis.SvcName}}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{.Redis.SvcName }}
spec:
  ports:
  - name: redis
    port: {{ .Redis.Port }}
  selector:
    app: {{ .Redis.SvcName }}
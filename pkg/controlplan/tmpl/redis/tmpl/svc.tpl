apiVersion: v1
kind: Service
metadata:
  name: {{.Spec.Redis.SvcName}}
  namespace: {{ ns . }}
  labels:
    app: {{.Spec.Redis.SvcName }}
spec:
  ports:
  - name: redis
    port: {{ .Spec.Redis.Port }}
  selector:
    app: {{ .Spec.Redis.SvcName }}
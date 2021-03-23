apiVersion: v1
kind: Service
metadata:
  name: {{.Spec.ControlPlan.Redis.SvcName}}
  namespace: {{ ns . }}
  labels:
    app: {{.Spec.ControlPlan.Redis.SvcName }}
spec:
  ports:
  - name: redis
    port: {{ .Spec.ControlPlan.Redis.Port }}
  selector:
    app: {{ .Spec.ControlPlan.Redis.SvcName }}
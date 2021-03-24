apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.ControlPlan.Pg.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: cnvrg-postgres
spec:
  ports:
    - port: {{ .Spec.ControlPlan.Pg.Port }}
  selector:
    app: {{ .Spec.ControlPlan.Pg.SvcName }}

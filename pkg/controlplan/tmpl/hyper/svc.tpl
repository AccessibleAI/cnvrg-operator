apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.ControlPlan.Hyper.SvcName }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Spec.ControlPlan.Hyper.SvcName }}
spec:
  ports:
    - port: {{ .Spec.ControlPlan.Hyper.Port }}
  selector:
    app: {{ .Spec.ControlPlan.Hyper.SvcName }}
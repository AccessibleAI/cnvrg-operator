apiVersion: v1
kind: Service
metadata:
  name: {{ .ControlPlan.Hyper.SvcName }}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{ .ControlPlan.Hyper.SvcName }}
spec:
  ports:
    - port: {{ .ControlPlan.Hyper.Port }}
  selector:
    app: {{ .ControlPlan.Hyper.SvcName }}
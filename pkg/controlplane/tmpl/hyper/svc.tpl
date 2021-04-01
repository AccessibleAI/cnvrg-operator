apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.ControlPlane.Hyper.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ .Spec.ControlPlane.Hyper.SvcName }}
spec:
  ports:
    - port: {{ .Spec.ControlPlane.Hyper.Port }}
  selector:
    app: {{ .Spec.ControlPlane.Hyper.SvcName }}
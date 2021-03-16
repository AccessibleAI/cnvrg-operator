apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Pg.SvcName }}
  namespace: {{ .Spec.CnvrgNs }}
  labels:
    app: cnvrg-postgres
spec:
  ports:
    - port: {{ .Spec.Pg.Port }}
  selector:
    app: {{ .Spec.Pg.SvcName }}

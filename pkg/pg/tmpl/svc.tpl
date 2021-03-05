apiVersion: v1
kind: Service
metadata:
  name: {{ .Pg.SvcName }}
  namespace: {{ .CnvrgNs }}
  labels:
    app: cnvrg-postgres
spec:
  ports:
    - port: {{.Pg.Port}}
  selector:
    app: {{ .Pg.SvcName }}

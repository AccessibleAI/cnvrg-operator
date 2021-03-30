apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Dbs.Pg.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: cnvrg-postgres
spec:
  ports:
    - port: {{ .Spec.Dbs.Pg.Port }}
  selector:
    app: {{ .Spec.Dbs.Pg.SvcName }}

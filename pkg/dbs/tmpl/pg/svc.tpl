apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Dbs.Pg.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: cnvrg-postgres
    owner: cnvrg-control-plane
spec:
  ports:
    - port: {{ .Spec.Dbs.Pg.Port }}
  selector:
    app: {{ .Spec.Dbs.Pg.SvcName }}
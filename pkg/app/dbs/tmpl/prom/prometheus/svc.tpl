apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Dbs.Prom.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
spec:
  ports:
    - name: http
      port: 9090
      targetPort: http
  selector:
    app: {{ .Spec.Dbs.Prom.SvcName }}
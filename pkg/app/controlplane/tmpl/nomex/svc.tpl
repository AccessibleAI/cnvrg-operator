apiVersion: v1
kind: Service
metadata:
  name: nomex
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
  labels:
    app: nomex
spec:
  ports:
  - port: 2112
  selector:
    app: nomex
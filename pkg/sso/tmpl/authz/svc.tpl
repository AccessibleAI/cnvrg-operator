apiVersion: v1
kind: Service
metadata:
  name: cnvrg-authz
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
spec:
  ports:
    - name: grpc
      port: 50052
  selector:
    app: cnvrg-authz
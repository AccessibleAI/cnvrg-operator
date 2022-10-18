apiVersion: v1
kind: Service
metadata:
  name: sso-central
  namespace: {{.Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
spec:
  ports:
    - name: http
      port: 8080
  selector:
    app: sso-central
apiVersion: v1
kind: Service
metadata:
  name: {{.Spec.SSO.Central.SvcName}}
  namespace: {{.Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
spec:
  ports:
    - name: http
      port: 8080
  selector:
    app: {{.Spec.SSO.Central.SvcName}}
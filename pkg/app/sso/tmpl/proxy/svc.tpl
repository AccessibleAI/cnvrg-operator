apiVersion: v1
kind: Service
metadata:
  name: {{.Spec.SSO.Proxy.SvcName}}
  namespace: {{.Namespace}}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
spec:
  ports:
    - name: http
      port: 80
      targetPort: 8888
  selector:
    app: {{.Spec.SSO.Proxy.SvcName}}
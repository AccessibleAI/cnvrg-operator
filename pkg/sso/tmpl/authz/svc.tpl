apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.SSO.Authz.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
spec:
  ports:
    - name: grpc
      port: 50052
  selector:
    app: {{ .Spec.SSO.Authz.SvcName }}
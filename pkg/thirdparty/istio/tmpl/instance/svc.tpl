apiVersion: v1
kind: Service
metadata:
  namespace:  {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
  labels:
    name: istio-operator
  name: istio-operator
spec:
  ports:
    - name: http-metrics
      port: 8383
      targetPort: 8383
  selector:
    name: istio-operator
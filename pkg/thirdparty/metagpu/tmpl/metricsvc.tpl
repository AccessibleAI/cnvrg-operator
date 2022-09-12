kind: Service
apiVersion: v1
metadata:
  name: metagpu-metrics-exporter
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
  labels:
    app: "metagpu-exporter"
spec:
  selector:
    name: metagpu-device-plugin
  ports:
    - protocol: TCP
      port: 2112
      name: metrics
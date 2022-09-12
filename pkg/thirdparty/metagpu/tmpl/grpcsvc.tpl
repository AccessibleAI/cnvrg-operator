kind: Service
apiVersion: v1
metadata:
  name: metagpu-device-plugin
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
  labels:
    app: "metagpu-device-plugin"
spec:
  selector:
    name: metagpu-device-plugin
  ports:
    - protocol: TCP
      port: 50052
      name: grcp
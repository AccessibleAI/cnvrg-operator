kind: Service
apiVersion: v1
metadata:
  name: metagpu-device-plugin
  namespace: {{ .Namespace }}
  labels:
    app: "metagpu-exporter"
spec:
  selector:
    name: metagpu-device-plugin
  ports:
    - protocol: TCP
      port: 50052
      name: grcp
    - protocol: TCP
      port: 2112
      name: metrics
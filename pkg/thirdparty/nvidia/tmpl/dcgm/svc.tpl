apiVersion: v1
kind: Service
metadata:
  name: dcgm-exporter
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
  labels:
    app: dcgm-exporter
spec:
  type: ClusterIP
  ports:
    - name: "metrics"
      port: 9400
      targetPort: 9400
      protocol: TCP
  selector:
    app: "dcgm-exporter"
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/name: node-exporter
    app.kubernetes.io/version: v1.0.1
    owner: cnvrg-control-plane
  name: node-exporter
  namespace: {{ ns . }}
spec:
  clusterIP: None
  ports:
  - name: https
    port: 9100
    targetPort: https
  selector:
    app.kubernetes.io/name: node-exporter

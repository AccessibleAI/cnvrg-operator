apiVersion: v1
kind: Service
metadata:
  name: dcgm-exporter
  namespace: {{ ns . }}
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
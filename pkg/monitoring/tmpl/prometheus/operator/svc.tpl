apiVersion: v1
kind: Service
metadata:
  name: cnvrg-prometheus-operator
  namespace: {{ ns . }}
  labels:
    app: cnvrg-prometheus-operator
    owner: cnvrg-control-plane
spec:
  ports:
  - name: https
    port: 8443
    targetPort: https
  selector:
    app: cnvrg-prometheus-operator

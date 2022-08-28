apiVersion: v1
kind: Service
metadata:
  name: prom
  namespace: {{ .Data.Namespace }}
spec:
  ports:
    - name: http
      port: 9090
  selector:
    app: prom
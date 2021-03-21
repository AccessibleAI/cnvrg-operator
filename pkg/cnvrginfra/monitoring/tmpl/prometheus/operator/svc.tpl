apiVersion: v1
kind: Service
metadata:
  name: cnvrg-prometheus-operator
  namespace: {{ .Spec.CnvrgInfraNs }}
  labels:
    app: cnvrg-prometheus-operator
spec:
  ports:
  - name: https
    port: 8443
    targetPort: https
  selector:
    app: cnvrg-prometheus-operator

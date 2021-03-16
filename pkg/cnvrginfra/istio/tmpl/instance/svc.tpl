apiVersion: v1
kind: Service
metadata:
  namespace:  {{ .Spec.CnvrgInfraNs }}
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
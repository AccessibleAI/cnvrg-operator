apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Prometheus.SvcName }}
  namespace: {{ .Namespace }}
  labels:
    app: cnvrg-ccp-prometheus
spec:
  sessionAffinity: ClientIP
  ports:
    - name: web
      port: {{ .Spec.Prometheus.Port }}
      targetPort: web
  selector:
    prometheus: cnvrg-ccp-prometheus


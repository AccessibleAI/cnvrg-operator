apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Monitoring.Prometheus.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: cnvrg-ccp-prometheus
spec:
  sessionAffinity: ClientIP
  ports:
    - name: web
      port: {{ .Spec.Monitoring.Prometheus.Port }}
      targetPort: web
  selector:
    prometheus: cnvrg-ccp-prometheus


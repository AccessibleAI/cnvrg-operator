apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Monitoring.Grafana.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: grafana
spec:
  ports:
  - name: http
    port: {{ .Spec.Monitoring.Grafana.Port }}
    targetPort: http
  selector:
    app: grafana

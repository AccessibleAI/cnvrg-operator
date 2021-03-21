apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Grafana.SvcName }}
  namespace: {{ .Namespace }}
  labels:
    app: grafana
spec:
  ports:
  - name: http
    port: {{ .Spec.Grafana.Port }}
    targetPort: http
  selector:
    app: grafana

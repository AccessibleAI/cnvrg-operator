apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Grafana.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: grafana
spec:
  ports:
  - name: http
    port: {{ .Spec.Grafana.Port }}
    targetPort: http
  selector:
    app: grafana

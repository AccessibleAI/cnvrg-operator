apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Monitoring.Prometheus.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: cnvrg-ccp-prometheus
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  sessionAffinity: ClientIP
  ports:
    - name: web
      port: {{ .Spec.Monitoring.Prometheus.Port }}
      targetPort: web
      {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
      nodePort: {{ .Spec.Monitoring.Prometheus.NodePort }}
      {{- end }}
  selector:
    app: prometheus


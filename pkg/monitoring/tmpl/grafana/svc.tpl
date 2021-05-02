apiVersion: v1
kind: Service
metadata:
  name: {{ .Spec.Monitoring.Grafana.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: grafana
    owner: cnvrg-control-plane
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  ports:
  - name: http
    port: {{ .Spec.Monitoring.Grafana.Port }}
    targetPort: http
    {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
    nodePort: {{ .Spec.Monitoring.Grafana.NodePort }}
    {{- end }}
  selector:
    app: grafana

apiVersion: v1
kind: Service
metadata:
  name: cnvrg-ingress-test
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: cnvrg-ingress-test
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
  type: NodePort
  {{- end }}
  sessionAffinity: ClientIP
  ports:
    - port: 8000
      targetPort: 8000
      {{- if eq .Spec.Networking.Ingress.Type "nodeport" }}
      nodePort: {{ .Spec.Monitoring.Prometheus.NodePort }}
      {{- end }}
  selector:
    app: ingresscheck


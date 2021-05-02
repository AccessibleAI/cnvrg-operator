apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Spec.Monitoring.Prometheus.SvcName }}
  namespace: {{ ns . }}
spec:
  rules:
    - host: "{{.Spec.Monitoring.Prometheus.SvcName}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            backend:
              service:
                name: {{ .Spec.Monitoring.Prometheus.SvcName }}
                port:
                  number: {{ .Spec.Monitoring.Prometheus.Port }}
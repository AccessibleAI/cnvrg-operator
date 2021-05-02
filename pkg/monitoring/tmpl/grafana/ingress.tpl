apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Spec.Monitoring.Grafana.SvcName }}
  namespace: {{ ns . }}
spec:
  rules:
    - host: "{{.Spec.Monitoring.Grafana.SvcName}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            backend:
              service:
                name: {{ .Spec.Monitoring.Grafana.SvcName }}
                port:
                  number: {{ .Spec.Monitoring.Grafana.Port }}
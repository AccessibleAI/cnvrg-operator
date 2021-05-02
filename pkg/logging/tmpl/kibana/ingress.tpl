apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
spec:
  rules:
    - host: "{{.Spec.Logging.Kibana.SvcName}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            backend:
              service:
                name: {{ .Spec.Logging.Kibana.SvcName }}
                port:
                  number: {{ .Spec.Logging.Kibana.Port }}

apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/proxy-send-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-read-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-body-size: 5G
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
spec:
  rules:
    - host: "{{.Spec.Logging.Kibana.SvcName}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: {{ .Spec.Logging.Kibana.SvcName }}
                port:
                  number: {{ .Spec.Logging.Kibana.Port }}

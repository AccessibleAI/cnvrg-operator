apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/proxy-send-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-read-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-body-size: 5G
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
spec:
  rules:
    - host: "{{.Spec.Logging.Kibana.SvcName}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            backend:
              serviceName: {{ .Spec.Logging.Kibana.SvcName }}
              servicePort: {{ .Spec.Logging.Kibana.SvcName }}

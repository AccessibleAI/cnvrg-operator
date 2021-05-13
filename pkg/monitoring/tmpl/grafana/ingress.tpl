apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/proxy-send-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-read-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-body-size: 5G
  name: {{ .Spec.Monitoring.Grafana.SvcName }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
spec:
  rules:
    - host: "{{.Spec.Monitoring.Grafana.SvcName}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            backend:
              serviceName: {{ .Spec.Monitoring.Grafana.SvcName }}
              servicePort: {{ .Spec.Monitoring.Grafana.Port }}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/proxy-send-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-read-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-body-size: 5G
    {{- range $k, $v := .Spec.Annotations }}
    {{ $k }}: "{{ $v }}"
    {{- end }}
  name: {{ .Spec.Monitoring.Prometheus.SvcName }}
  namespace: {{ ns . }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{ $k }}: "{{ $v }}"
    {{- end }}
spec:
  {{- if ne .Spec.Networking.HTTPS.CertSecret "" }}
  tls:
  - hosts:
    - "{{ .Spec.Monitoring.Prometheus.SvcName}}.{{ .Spec.ClusterDomain }}"
    secretName: "{{ .Spec.Networking.HTTPS.CertSecret }}"
  {{- end }}
  rules:
  - host: "{{ .Spec.Monitoring.Prometheus.SvcName }}.{{ .Spec.ClusterDomain }}"
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: {{ .Spec.Monitoring.Prometheus.SvcName }}
            port:
              number: {{ .Spec.Monitoring.Prometheus.Port }}

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
    {{- if and ( isTrue .Spec.Networking.Ingress.DynamicCertsEnabled ) (ne .Spec.Networking.Ingress.DynamicCertsIssuer "") }}
    cert-manager.io/issuer: {{ .Spec.Networking.Ingress.DynamicCertsIssuer }}
    {{- end }}
  name: {{ .Spec.Monitoring.Prometheus.SvcName }}
  namespace: {{ ns . }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{ $k }}: "{{ $v }}"
    {{- end }}
spec:
  {{- if isTrue .Spec.Networking.HTTPS.Enabled }}
  tls:
  - hosts:
      - {{ .Spec.Monitoring.Prometheus.SvcName}}.{{ .Spec.ClusterDomain }}
    {{- if ne .Spec.Networking.HTTPS.CertSecret "" }}
    secretName: {{ .Spec.Networking.HTTPS.CertSecret }}
    {{- end }}
    {{- if isTrue .Spec.Networking.Ingress.DynamicCertsEnabled }}
    secretName: {{ .Spec.Monitoring.Prometheus.SvcName }}-tls
    {{- end }}
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

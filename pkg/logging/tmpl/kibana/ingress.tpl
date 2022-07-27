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
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{ $k }}: "{{ $v }}"
    {{- end }}
spec:
{{- if ne .Spec.Networking.Ingress.ClassName "" }}
  ingressClassName: {{ .Spec.Networking.Ingress.ClassName }}
{{- end }}
  rules:
  - host: "{{ .Spec.Logging.Kibana.SvcName }}.{{ .Spec.ClusterDomain }}"
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: {{ .Spec.Logging.Kibana.SvcName }}
            port:
              number: {{ .Spec.Logging.Kibana.Port }}
{{- if isTrue .Spec.Networking.HTTPS.Enabled }}
  tls:
  - hosts:
    - {{ .Spec.Logging.Kibana.SvcName }}.{{ .Spec.ClusterDomain }}
{{- if isTrue .Spec.Networking.HTTPS.AcmeCert }}
    secretName: kibana-tls-cert
{{- else if ne .Spec.Networking.HTTPS.CertSecret "" }}
    secretName: {{ .Spec.Networking.HTTPS.CertSecret }}
{{- end }}
{{- end }}

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
  name: {{ .Spec.Dbs.Es.SvcName }}
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
  - host: "{{ .Spec.Dbs.Es.SvcName }}.{{ .Spec.ClusterDomain }}"
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: {{ .Spec.Dbs.Es.SvcName }}
            port:
              number: {{ .Spec.Dbs.Es.Port }}
{{- if isTrue .Spec.Networking.HTTPS.Enabled }}
  tls:
  - hosts:
    - {{ .Spec.Dbs.Es.SvcName }}.{{ .Spec.ClusterDomain }}
{{- if isTrue .Spec.Networking.HTTPS.AcmeCert }}
    secretName: es-tls-cert
{{- else if ne .Spec.Networking.HTTPS.CertSecret "" }}
    secretName: {{ .Spec.Networking.HTTPS.CertSecret }}
{{- end }}
{{- end }}

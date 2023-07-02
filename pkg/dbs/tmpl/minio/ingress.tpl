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
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{ ns . }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{ $k }}: "{{ $v }}"
    {{- end }}
spec:
  {{- if isTrue .Spec.Networking.HTTPS.Enabled }}
  tls:
  - hosts:
      - {{ .Spec.Dbs.Minio.SvcName}}.{{ .Spec.ClusterDomain }}
    {{- if ne .Spec.Networking.HTTPS.CertSecret "" }}
    secretName: {{ .Spec.Networking.HTTPS.CertSecret }}
    {{- end }}
    {{- if isTrue .Spec.Networking.Ingress.DynamicCertsEnabled }}
    secretName: {{ .Spec.Dbs.Minio.SvcName }}-tls
    {{- end }}
  {{- end }}
  rules:
  - host: "{{ .Spec.Dbs.Minio.SvcName }}.{{ .Spec.ClusterDomain }}"
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: {{ .Spec.Dbs.Minio.SvcName }}
            port:
              number: {{ .Spec.Dbs.Minio.Port }}

apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    nginx.ingress.kubernetes.io/proxy-send-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-read-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-body-size: 5G
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Spec.SSO.Jwks.SvcName }}
  namespace: {{.Namespace }}
spec:
  ingressClassName: nginx
  {{- if and ( isTrue .Spec.Networking.HTTPS.Enabled ) (ne .Spec.Networking.HTTPS.CertSecret "") }}
  tls:
  - hosts:
      - {{.Spec.SSO.Jwks.SvcName}}{{.Spec.Networking.ClusterDomainPrefix.Prefix}}.{{ .Spec.ClusterDomain }}
    secretName: {{ .Spec.Networking.HTTPS.CertSecret }}
  {{- end }}
  rules:
    - host: "{{.Spec.SSO.Jwks.SvcName}}{{.Spec.Networking.ClusterDomainPrefix.Prefix}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: "{{.Spec.SSO.Jwks.SvcName}}"
                port:
                  number: 8080
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/proxy-send-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-read-timeout: 18000s
    nginx.ingress.kubernetes.io/proxy-body-size: 5G
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- if isTrue .Spec.SSO.Enabled }}
    sso.cnvrg.io/enabled: "true"
    sso.cnvrg.io/skipAuthRoutes: \/api\/health
    sso.cnvrg.io/central: "{{ .Spec.SSO.Central.PublicUrl }}"
    sso.cnvrg.io/upstream: "{{.Spec.Dbs.Prom.Grafana.SvcName}}.{{ .Namespace }}.svc:{{.Spec.Dbs.Prom.Grafana.Port}}"
    {{- end }}
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Spec.Dbs.Prom.Grafana.SvcName }}
  namespace: {{ .Namespace }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  ingressClassName: nginx
  {{- if and ( isTrue .Spec.Networking.HTTPS.Enabled ) (ne .Spec.Networking.HTTPS.CertSecret "") }}
  tls:
  - hosts:
      - {{ .Spec.Monitoring.Grafana.SvcName}}.{{ .Spec.ClusterDomain }}
    secretName: {{ .Spec.Networking.HTTPS.CertSecret }}
  {{- end }}
  rules:
    - host: "{{.Spec.Dbs.Prom.Grafana.SvcName}}.{{ .Spec.ClusterDomain }}"
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                {{- if isTrue .Spec.SSO.Enabled }}
                name: {{ .Spec.SSO.Proxy.SvcName }}
                port:
                  number: 80
                {{- else }}
                name: {{ .Spec.Dbs.Prom.Grafana.SvcName }}
                port:
                  number: {{ .Spec.Dbs.Prom.Grafana.Port }}
                {{- end }}

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
    sso.cnvrg.io/skipAuthRoutes: \/assets \/healthz \/public \/pack \/vscode.tar.gz \/jupyter.vsix \/gitlens.vsix \/ms-python-release.vsix \/webhooks \/api/v2/metrics \/api/v1/events/endpoint_rule_alert \/api/v2/version
    sso.cnvrg.io/central: "{{ .Spec.SSO.Central.PublicUrl }}"
    sso.cnvrg.io/upstream: "{{ .Spec.ControlPlane.WebApp.SvcName }}.{{ .Namespace }}.svc:{{.Spec.ControlPlane.WebApp.Port}}"
    {{- end }}
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
  namespace: {{ .Namespace }}
spec:
  ingressClassName: nginx
  {{- if and ( isTrue .Spec.Networking.HTTPS.Enabled ) (ne .Spec.Networking.HTTPS.CertSecret "") }}
  tls:
  - hosts:
      - {{ .Spec.ControlPlane.WebApp.SvcName}}{{.Spec.Networking.ClusterDomainPrefix.Prefix}}.{{ .Spec.ClusterDomain }}
    secretName: {{ .Spec.Networking.HTTPS.CertSecret }}
  {{- end }}
  rules:
    - host: "{{.Spec.ControlPlane.WebApp.SvcName}}{{.Spec.Networking.ClusterDomainPrefix.Prefix}}.{{ .Spec.ClusterDomain }}"
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
                name: {{ .Spec.ControlPlane.WebApp.SvcName }}
                port:
                  number: {{ .Spec.ControlPlane.WebApp.Port }}
                {{- end}}

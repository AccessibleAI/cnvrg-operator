apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
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
spec:
  hosts:
    - "{{.Spec.ControlPlane.WebApp.SvcName}}{{.Spec.Networking.ClusterDomainPrefix.Prefix}}.{{ .Spec.ClusterDomain }}"
  gateways:
    - {{ .Spec.Networking.Ingress.IstioGwName}}
  http:
    - retries:
        attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
        perTryTimeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
      timeout: {{ .Spec.Networking.Ingress.Timeout }}
      route:
        - destination:
            {{- if isTrue .Spec.SSO.Enabled }}
            host: "{{.Spec.SSO.Proxy.Address}}"
            {{- else }}
            host: "{{ .Spec.ControlPlane.WebApp.SvcName }}.{{ .Namespace }}.svc.{{ .Spec.ClusterInternalDomain }}"
            {{- end }}
      headers:
        request:
          set:
            {{- if isTrue .Spec.Networking.HTTPS.Enabled }}
            x-forwarded-proto: https
            {{- else }}
            x-forwarded-proto: http
            {{- end}}


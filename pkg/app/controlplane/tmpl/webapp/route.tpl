apiVersion: route.openshift.io/v1
kind: Route
metadata:
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    haproxy.router.openshift.io/timeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
    {{- if isTrue .Spec.SSO.Enabled }}
    sso.cnvrg.io/enabled: "true"
    sso.cnvrg.io/skipAuthRoutes: \/assets \/healthz \/public \/pack \/vscode.tar.gz \/jupyter.vsix \/gitlens.vsix \/ms-python-release.vsix \/webhooks \/api/v2/metrics \/api/v1/events/endpoint_rule_alert \/api/v2/version
    sso.cnvrg.io/central: "{{ .Spec.SSO.Central.PublicUrl }}"
    sso.cnvrg.io/upstream: "{{ .Spec.ControlPlane.WebApp.SvcName }}.{{ .Namespace }}.svc:{{.Spec.ControlPlane.WebApp.Port}}"
    {{- end }}
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Namespace }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  host: "{{ .Spec.ControlPlane.WebApp.SvcName }}{{.Spec.Networking.ClusterDomainPrefix.Prefix}}.{{ .Spec.ClusterDomain }}"
  port:
    {{- if isTrue .Spec.SSO.Enabled }}
    targetPort: 8888
    {{- else }}
    targetPort: {{ .Spec.ControlPlane.WebApp.Port }}
    {{- end }}
  to:
    kind: Service
    {{- if isTrue .Spec.SSO.Enabled }}
    name: {{ .Spec.SSO.Proxy.SvcName }}
    weight: 100
    {{- else }}
    name: {{ .Spec.ControlPlane.WebApp.SvcName }}
    weight: 100
    {{- end}}
  {{- if isTrue .Spec.Networking.HTTPS.Enabled  }}
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
    {{- if isTrue .Spec.Networking.Ingress.OcpSecureRoutes }}
    certificate: |
{{ printf "%s" .Spec.Networking.HTTPS.Cert | indent 6 }}
    key: |
{{ printf "%s" .Spec.Networking.HTTPS.Key | indent 6 }}
    {{- end }}
  {{- end }}
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    haproxy.router.openshift.io/timeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: cnvrg-jwks
  namespace: {{.Namespace }}
  labels:
    app: {{.Namespace }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  host: "{{.Spec.SSO.Jwks.Name}}.{{ .Spec.ClusterDomain }}"
  port:
    targetPort: 8080
  to:
    kind: Service
    name: "{{.Spec.SSO.Jwks.Name}}.{{.Namespace }}.svc"
    weight: 100
  {{- if isTrue .Spec.Networking.HTTPS.Enabled  }}
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
  {{- end }}
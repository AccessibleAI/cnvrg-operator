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
  name: {{.Spec.SSO.Central.SvcName}}
  namespace: {{.Namespace }}
  labels:
    app: {{.Namespace }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  host: "{{.Spec.SSO.Central.SvcName}}.{{ .Spec.ClusterDomain }}"
  port:
    targetPort: 8080
  to:
    kind: Service
    name: "{{.Spec.SSO.Central.SvcName}}"
    weight: 100
  {{- if isTrue .Spec.Networking.HTTPS.Enabled  }}
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
    ### secure route section placeholder start ###
    {{- if and ( isTrue .Spec.Networking.Ingress.OcpSecureRoutes ) (ne .Spec.Networking.HTTPS.CertSecret "") }}
    certificate: tls_crt_content
    key: tls_key_content
    {{- end }}
    ### secure route section placeholder end ###
  {{- end }}
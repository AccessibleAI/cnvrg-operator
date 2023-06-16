apiVersion: route.openshift.io/v1
kind: Route
metadata:
  annotations:
    haproxy.router.openshift.io/timeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
  name: {{ .Spec.Monitoring.Prometheus.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ ns . }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  host: "{{ .Spec.Monitoring.Prometheus.SvcName }}.{{ .Spec.ClusterDomain }}"
  port:
    targetPort: {{ .Spec.Monitoring.Prometheus.Port }}
  to:
    kind: Service
    name: {{ .Spec.Monitoring.Prometheus.SvcName }}
    weight: 100
  {{- if isTrue .Spec.Networking.HTTPS.Enabled  }}
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
    {{- if isTrue .Spec.Networking.Ingress.OcpSecureRoutes }}
    {{- $tlsSecret := secret .Spec.Networking.HTTPS.CertSecret }}
    certificate: |-
      {{ $tlsSecret.Data."tls.crt" | indent 6 }}
    key: |-
      {{ $tlsSecret.Data."tls.key" | indent 6 }}
    {{- end }}
  {{- end }}
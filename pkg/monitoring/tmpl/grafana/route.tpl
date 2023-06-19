apiVersion: route.openshift.io/v1
kind: Route
metadata:
  annotations:
    haproxy.router.openshift.io/timeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Spec.Monitoring.Grafana.SvcName }}
  namespace: {{ ns . }}

  labels:
    app: {{ ns . }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  host: "{{ .Spec.Monitoring.Grafana.SvcName }}.{{ .Spec.ClusterDomain }}"
  port:
    targetPort: {{ .Spec.Monitoring.Grafana.Port }}
  to:
    kind: Service
    name: {{ .Spec.Monitoring.Grafana.SvcName }}
    weight: 100
  {{- if isTrue .Spec.Networking.HTTPS.Enabled  }}
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
    {{- if isTrue .Spec.Networking.Ingress.OcpSecureRoutes }}
    certificate: |-
      {{- $tlsCrtVal := (lookup "v1" "Secret" "cnvrg" (print .Spec.Networking.HTTPS.CertSecret)).data.tls\.crt | default dict }}
      {{ $tlsCrtVal | indent 6 }}
    key: |-
      {{- $tlsKeyVal := (lookup "v1" "Secret" "cnvrg" (print .Spec.Networking.HTTPS.CertSecret)).data.tls\.key | default dict }}
      {{ $tlsKeyVal | indent 6 }}
    {{- end }}
  {{- end }}

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
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{.Namespace }}
  labels:
    app: {{.Namespace }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  host: "{{ .Spec.Dbs.Minio.SvcName }}.{{ .Spec.ClusterDomain }}"
  port:
    targetPort: {{ .Spec.Dbs.Minio.Port }}
  to:
    kind: Service
    name: {{ .Spec.Dbs.Minio.SvcName }}
    weight: 100
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
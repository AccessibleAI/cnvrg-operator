{{- import "sigs.k8s.io/controller-runtime/pkg/template" }}
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  annotations:
    haproxy.router.openshift.io/timeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ ns . }}
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
    certificate: |-
      {{- $tlsCrtVal := (template.lookup "v1" "Secret" "cnvrg" (print .Spec.Networking.HTTPS.CertSecret)).data.tls\.crt | default dict }}
      {{ $tlsCrtVal | indent 6 }}
    key: |-
      {{- $tlsKeyVal := (template.lookup "v1" "Secret" "cnvrg" (print .Spec.Networking.HTTPS.CertSecret)).data.tls\.key | default dict }}
      {{ $tlsKeyVal | indent 6 }}
    {{- end }}
  {{- end }}

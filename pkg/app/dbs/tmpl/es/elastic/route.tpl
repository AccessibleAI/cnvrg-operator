apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name: {{ .Spec.Dbs.Es.SvcName }}
  namespace: {{ ns . }}
  annotations:
    haproxy.router.openshift.io/timeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    app: {{ ns . }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  host: "{{ .Spec.Dbs.Es.SvcName }}.{{ .Spec.ClusterDomain }}"
  port:
    targetPort: {{ .Spec.Dbs.Es.Port }}
  to:
    kind: Service
    name: {{ .Spec.Dbs.Es.SvcName }}
    weight: 100
  {{- if isTrue .Spec.Networking.HTTPS.Enabled  }}
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
  {{- end }}
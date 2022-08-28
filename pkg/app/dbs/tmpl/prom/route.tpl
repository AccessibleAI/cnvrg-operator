apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name: prom
  namespace: {{ .Data.Namespace }}
  labels:
    app: {{ ns . }}
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  host: "prom.{{ .Data.ClusterDomain }}"
  port:
    targetPort: 9090
  to:
    kind: Service
    name: prom
    weight: 100
  {{- if isTrue .Data.HttpsEnabled  }}
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
  {{- end }}
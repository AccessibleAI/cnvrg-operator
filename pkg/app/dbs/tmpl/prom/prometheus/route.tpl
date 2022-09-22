apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name: prom
  namespace: {{ .Namespace }}
  labels:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    app: {{ ns . }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  host: "prom.{{ .Spec.ClusterDomain }}"
  port:
    targetPort: 9090
  to:
    kind: Service
    name: prom
    weight: 100
  {{- if isTrue .Spec.Networking.HTTPS.Enabled  }}
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
  {{- end }}
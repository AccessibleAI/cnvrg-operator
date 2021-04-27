apiVersion: route.openshift.io/v1
kind: Route
metadata:
  annotations:
    haproxy.router.openshift.io/timeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
  name: {{ .Spec.Dbs.Es.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ ns . }}
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
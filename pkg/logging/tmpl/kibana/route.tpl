apiVersion: route.openshift.io/v1
kind: Route
metadata:
  annotations:
    haproxy.router.openshift.io/timeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
  labels:
    app: {{ ns . }}
spec:
  host: "{{ .Spec.Logging.Kibana.SvcName }}.{{ .Spec.ClusterDomain }}"
  port:
    targetPort: {{ .Spec.Logging.Kibana.Port }}
  to:
    kind: Service
    name: {{ .Spec.Logging.Kibana.SvcName }}
    weight: 100
  {{- if isTrue .Spec.Networking.HTTPS.Enabled  }}
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
  {{- end }}
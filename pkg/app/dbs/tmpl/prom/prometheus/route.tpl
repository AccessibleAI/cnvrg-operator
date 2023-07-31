apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name:  {{ .Spec.Dbs.Prom.SvcName }}
  namespace: {{ .Namespace }}
  labels:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    app: {{ .Namespace }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  host: "{{ .Spec.Dbs.Prom.SvcName }}.{{ .Spec.ClusterDomain }}"
  port:
    targetPort: {{ .Spec.Dbs.Prom.Port }}
  to:
    kind: Service
    name:  {{ .Spec.Dbs.Prom.SvcName }}
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
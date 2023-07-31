apiVersion: route.openshift.io/v1
kind: Route
metadata:
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    haproxy.router.openshift.io/timeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
    {{- if isTrue .Spec.SSO.Enabled }}
    sso.cnvrg.io/enabled: "true"
    sso.cnvrg.io/skipAuthRoutes: ""
    sso.cnvrg.io/central: "{{ .Spec.SSO.Central.PublicUrl }}"
    sso.cnvrg.io/upstream: "{{ .Spec.Dbs.Es.Kibana.SvcName }}.{{ .Namespace }}.svc:{{.Spec.Dbs.Es.Kibana.Port}}"
    {{- end }}
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  name: {{ .Spec.Dbs.Es.Kibana.SvcName }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Namespace }}
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  host: "{{ .Spec.Dbs.Es.Kibana.SvcName }}.{{ .Spec.ClusterDomain }}"
  port:
    {{- if isTrue .Spec.SSO.Enabled }}
    targetPort: 8888
    {{- else }}
    targetPort: {{ .Spec.Dbs.Es.Kibana.Port }}
    {{- end }}
  to:
    kind: Service
    {{- if isTrue .Spec.SSO.Enabled }}
    name: {{ .Spec.SSO.Proxy.SvcName }}
    weight: 100
    {{- else }}
    name: {{ .Spec.Dbs.Es.Kibana.SvcName }}
    weight: 100
    {{- end}}
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
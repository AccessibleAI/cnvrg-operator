apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.Dbs.Es.Kibana.SvcName }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- if isTrue .Spec.SSO.Enabled }}
    sso.cnvrg.io/enabled: "true"
    sso.cnvrg.io/skipAuthRoutes: ""
    sso.cnvrg.io/central: "{{ .Spec.SSO.Central.PublicUrl }}"
    sso.cnvrg.io/upstream: "{{ .Spec.Dbs.Es.Kibana.SvcName }}.{{ ns . }}.svc:{{.Spec.Dbs.Es.Kibana.Port}}"
    {{- end }}
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  hosts:
    - "{{.Spec.Dbs.Es.Kibana.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
  - {{ .Spec.Networking.Ingress.IstioGwName}}
  http:
  - retries:
      attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
      perTryTimeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
    timeout: {{ .Spec.Networking.Ingress.Timeout }}
    route:
    - destination:
        port:
          number: {{.Spec.Dbs.Es.Kibana.Port}}
        host: "{{ .Spec.Dbs.Es.Kibana.SvcName }}.{{ .Namespace }}.svc.{{ .Spec.ClusterInternalDomain }}"
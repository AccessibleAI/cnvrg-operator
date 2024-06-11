apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Spec.SSO.Jwks.SvcName }}
  namespace: {{.Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  config.yaml: |-
    discovery:
      secret:
        namespace: {{.Namespace}}
        labelKey: domainId
        dataKey: CNVRG_PKI_PUBLIC_KEY
    cache:
      enabled: true
      redis:
        address: localhost:6379
    api:
      listen: 0.0.0.0:8080
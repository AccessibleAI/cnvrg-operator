apiVersion: v1
kind: Secret
metadata:
  name: {{ .Data.Pki.PublicKeySecret }}
  namespace: {{ .Namespace }}
  annotations:
  {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
  {{- end }}
  labels:
    domainId: {{ .Data.DomainID }}
  {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
  {{- end }}
data:
  CNVRG_PKI_PUBLIC_KEY: {{ .Data.Keys.PublicKey | b64enc }}

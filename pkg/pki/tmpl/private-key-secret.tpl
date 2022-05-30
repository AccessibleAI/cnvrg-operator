apiVersion: v1
kind: Secret
metadata:
  name: {{ .Data.Pki.PrivateKeySecret }}
  namespace: {{ .Namespace }}
  annotations:
  {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
  {{- end }}
  labels:
  {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
  {{- end }}
data:
  CNVRG_PKI_PRIVATE_KEY: {{ .Data.Keys.PrivateKey | b64enc }}

{{- $altNames := list ( printf "%s.%s" .Spec.Dbs.Es.SvcName (.Namespace)) ( printf "%s.%s.svc" .Spec.Dbs.Es.SvcName (.Namespace) ) -}}
{{- $ca := genCA "elasticsearch-ca" 3650 -}}
{{- $cert := genSignedCert .Spec.Dbs.Es.SvcName nil $altNames 3650 $ca -}}
apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.Dbs.Es.SvcName }}-certs
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "false"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  tls.crt: {{ $cert.Cert | toString | b64enc }}
  tls.key: {{ $cert.Key | toString | b64enc }}
  ca.crt: {{ $ca.Cert | toString | b64enc }}
{{- $altNames := list ( printf "%s.%s" .Spec.Dbs.Es.SvcName (ns .)) ( printf "%s.%s.svc" .Spec.Dbs.Es.SvcName (ns .) ) -}}
{{- $ca := genCA "elasticsearch-ca" 3650 -}}
{{- $cert := genSignedCert .Spec.Dbs.Es.SvcName nil $altNames 3650 $ca -}}
apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.Dbs.Es.ServiceAccount }}-certs
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  tls.crt: {{ $cert.Cert | toString | b64enc }}
  tls.key: {{ $cert.Key | toString | b64enc }}
  ca.crt: {{ $ca.Cert | toString | b64enc }}
apiVersion: v1
kind: Secret
metadata:
  name: cvat-pg-secret
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
  CNVRG_CVAT_POSTGRES_PASSWORD: {{ randAlphaNum 128 | lower | b64enc }}

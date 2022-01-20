apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.Dbs.Cvat.Redis.CredsRef }}
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
  CNVRG_CVAT_REDIS_HOST: {{ .Spec.Dbs.Cvat.Redis.SvcName | b64enc }}

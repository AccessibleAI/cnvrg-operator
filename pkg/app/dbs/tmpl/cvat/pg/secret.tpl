apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.Dbs.Cvat.Pg.CredsRef }}
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
  {{- $user := "cnvrg" | b64enc}}
  {{- $pass := randAlphaNum 20 | b64enc }}
  {{- $db := "cnvrg_cvat" | b64enc }}

  # required vars for the PG-DB (could be omitted when when using external PG instance)
  POSTGRESQL_DATABASE: {{ $db }}
  POSTGRESQL_USER: {{ $user }}
  POSTGRESQL_PASSWORD: {{ $pass }}
  POSTGRESQL_ADMIN_PASSWORD: {{ $pass }}

  # required vars for the app
  CNVRG_CVAT_POSTGRES_DBNAME: {{ $db }}
  CNVRG_CVAT_POSTGRES_USER: {{ $user }}
  CNVRG_CVAT_POSTGRES_PASSWORD: {{ $pass }}
  CNVRG_CVAT_POSTGRES_HOST: {{ .Spec.Dbs.Cvat.Pg.SvcName | b64enc }}

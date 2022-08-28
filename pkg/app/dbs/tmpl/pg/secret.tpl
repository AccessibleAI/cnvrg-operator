apiVersion: v1
kind: Secret
metadata:
  name: {{ .Data.CredsRef }}
  namespace: {{ .Data.Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  {{- $user := "cnvrg" | b64enc}}
  {{- $pass := randAlphaNum 20 | b64enc }}
  {{- $db := "cnvrg_production" | b64enc }}
  # required vars for the PG-DB (could be omitted when when using external PG instance)
  POSTGRESQL_USER:                  {{ $user }}
  POSTGRESQL_PASSWORD:              {{ $pass }}
  POSTGRESQL_ADMIN_PASSWORD:        {{ $pass }}
  POSTGRESQL_DATABASE:              {{ $db }}
  POSTGRESQL_MAX_CONNECTIONS:       {{ .Data.MaxConnections | toString | b64enc }}
  POSTGRESQL_SHARED_BUFFERS:        {{ .Data.SharedBuffers | b64enc }}
  POSTGRESQL_EFFECTIVE_CACHE_SIZE:  {{ .Data.EffectiveCacheSize | b64enc }}
  # required vars for the app
  POSTGRES_DB:                      {{ $db }}
  POSTGRES_PASSWORD:                {{ $pass }}
  POSTGRES_USER:                    {{ $user }}
  POSTGRES_HOST:                    {{ .Data.SvcName | b64enc }}
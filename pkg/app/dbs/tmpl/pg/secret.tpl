apiVersion: v1
kind: Secret
metadata:
  name: {{ .CredsRef }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "false"
    {{- range $k, $v := .Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    {{- range $k, $v := .Labels }}
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
  POSTGRESQL_MAX_CONNECTIONS:       {{ .MaxConnections | toString | b64enc }}
  POSTGRESQL_SHARED_BUFFERS:        {{ .SharedBuffers | b64enc }}
  POSTGRESQL_EFFECTIVE_CACHE_SIZE:  {{ .EffectiveCacheSize | b64enc }}
  # required vars for the app
  POSTGRES_DB:                      {{ $db }}
  POSTGRES_PASSWORD:                {{ $pass }}
  POSTGRES_USER:                    {{ $user }}
  POSTGRES_HOST:                    {{ .SvcName | b64enc }}
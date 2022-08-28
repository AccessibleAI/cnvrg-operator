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
  {{- $pass := randAlphaNum 20 | b64enc }}
  {{- $redisUrl := printf ("redis://:%s@%s") $pass .Data.SvcName | b64enc }}
  {{- $conf := redisConf $pass | b64enc }}
  CNVRG_REDIS_PASSWORD:               {{ $pass }}
  OAUTH2_PROXY_REDIS_CONNECTION_URL:  {{ $redisUrl }} # for oauth2 proxy
  REDIS_URL:                          {{ $redisUrl }} # for cnvrg webapp/sidekiq
  redis.conf:                         {{ $conf }}
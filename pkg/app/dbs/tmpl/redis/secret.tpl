apiVersion: v1
kind: Secret
metadata:
  name: {{ .Spec.Dbs.Redis.CredsRef }}
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
  {{- $pass := randAlphaNum 20 | b64enc }}
  {{- $redisUrl := printf ("redis://:%s@%s") $pass .Spec.Dbs.Redis.SvcName | b64enc }}
  {{- $conf := redisConf $pass | b64enc }}
  CNVRG_REDIS_PASSWORD:               {{ $pass }}
  OAUTH2_PROXY_REDIS_CONNECTION_URL:  {{ $redisUrl }} # for oauth2 proxy
  REDIS_URL:                          {{ $redisUrl }} # for cnvrg webapp/sidekiq
  redis.conf:                         {{ $conf }}
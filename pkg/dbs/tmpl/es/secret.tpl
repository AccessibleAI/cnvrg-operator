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
  {{- $esUser := "cnvrg" | b64enc}}
  {{- $esPass := randAlphaNum 20 | b64enc}}
  {{- $esUrl  := printf ("http://%s:%s@%s") $esUser $esPass .Data.EsUrl | b64enc }} # env for webapp/kiqs
  CNVRG_ES_USER:          {{ $esUser }}     # env for webapp/kiqs
  CNVRG_ES_PASS:          {{ $esPass }}     # env for webapp/kiqs
  ELASTICSEARCH_URL:      {{ $esUrl  }}     # env for webapp/kiqs
  ES_USERNAME:            {{ $esUser }}     # env for elastalerts
  ES_PASSWORD:            {{ $esPass }}     # env for elastalerts
  ELASTICSEARCH_USERNAME: {{ $esUser }}     # env for kibana
  ELASTICSEARCH_PASSWORD: {{ $esPass }}     # env for kibana
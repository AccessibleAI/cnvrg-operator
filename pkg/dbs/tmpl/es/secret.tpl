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
  {{- $esUser := "cnvrg" }}
  {{- $esPass := randAlphaNum 20 }}
  {{- $esUrl  := printf ("http://%s:%s@%s") $esUser $esPass .Data.EsUrl }} # env for webapp/kiqs
  CNVRG_ES_USER:          {{ $esUser | b64enc }}     # env for webapp/kiqs
  CNVRG_ES_PASS:          {{ $esPass | b64enc }}     # env for webapp/kiqs
  ELASTICSEARCH_URL:      {{ $esUrl  | b64enc }}     # env for webapp/kiqs
  ES_USERNAME:            {{ $esUser | b64enc }}     # env for elastalerts
  ES_PASSWORD:            {{ $esPass | b64enc }}     # env for elastalerts
  ELASTICSEARCH_USERNAME: {{ $esUser | b64enc }}     # env for kibana
  ELASTICSEARCH_PASSWORD: {{ $esPass | b64enc }}     # env for kibana
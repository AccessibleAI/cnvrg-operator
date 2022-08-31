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
  {{- $esUser := "cnvrg" }}
  {{- $esPass := randAlphaNum 20 }}
  {{- $esUrl  := printf ("http://%s:%s@%s") $esUser $esPass .EsUrl }} # env for webapp/kiqs
  CNVRG_ES_USER:          {{ $esUser | b64enc }}     # env for webapp/kiqs
  CNVRG_ES_PASS:          {{ $esPass | b64enc }}     # env for webapp/kiqs
  ELASTICSEARCH_URL:      {{ $esUrl  | b64enc }}     # env for webapp/kiqs
  ES_USERNAME:            {{ $esUser | b64enc }}     # env for elastalerts
  ES_PASSWORD:            {{ $esPass | b64enc }}     # env for elastalerts
  ELASTICSEARCH_USERNAME: {{ $esUser | b64enc }}     # env for kibana
  ELASTICSEARCH_PASSWORD: {{ $esPass | b64enc }}     # env for kibana
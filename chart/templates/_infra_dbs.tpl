{{- define "spec.infra_dbs" }}
dbs:
  redis:
    enabled: {{ .Values.dbs.redis.enabled }}
    storageSize: {{ .Values.dbs.redis.storageSize }}
    storageClass: "{{ .Values.dbs.redis.storageClass }}"
    nodeSelector:
    {{- range $key, $value := .Values.dbs.redis.nodeSelector }}
      {{$key}}: {{$value}}
    {{- end }}
{{- end }}
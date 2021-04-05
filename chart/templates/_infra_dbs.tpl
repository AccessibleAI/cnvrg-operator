{{- define "spec.infra_dbs" }}
dbs:
  redis:
    {{- if eq (.Values.namespaceTenancy|toString) "true" }}
    enabled: "{{ .Values.dbs.redis.enabled }}"
    {{- else }}
    enabled: "false"
    {{- end }}
    appendonly: "{{ .Values.dbs.redis.appendonly }}"
    image: {{ .Values.dbs.redis.image }}
    limits:
      cpu: {{ .Values.dbs.redis.limits.cpu }}
      memory: {{ .Values.dbs.redis.limits.memory }}
    port: {{ .Values.dbs.redis.port }}
    requests:
      cpu: {{ .Values.dbs.redis.requests.cpu }}
      memory: {{ .Values.dbs.redis.requests.memory }}
    serviceAccount: {{ .Values.dbs.redis.serviceAccount }}
    storageSize: {{ .Values.dbs.redis.storageSize }}
    svcName: {{ .Values.dbs.redis.svcName }}
{{- end }}
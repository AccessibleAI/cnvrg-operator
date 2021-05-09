{{- define "spec.app_dbs" }}
dbs:
  es:
    enabled: {{ .Values.dbs.es.enabled }}
    storageSize: {{ .Values.dbs.es.storageSize }}
    storageClass: {{ .Values.dbs.es.storageClass }}
    nodeSelector:
    {{- range $key, $value := .Values.dbs.es.nodeSelector }}
      - {{$key}}: {{$value}}
    {{- end }}

  minio:
    enabled: {{ .Values.dbs.minio.enabled }}
    storageSize: {{ .Values.dbs.minio.storageSize }}
    storageClass: {{ .Values.dbs.minio.storageClass }}
    nodeSelector:
    {{- range $key, $value := .Values.dbs.minio.nodeSelector }}
      - {{$key}}: {{$value}}
    {{- end }}

  pg:
    enabled: {{ .Values.dbs.pg.enabled }}
    storageSize: {{ .Values.dbs.pg.storageSize }}
    storageClass: {{ .Values.dbs.pg.storageClass }}
    nodeSelector:
    {{- range $key, $value := .Values.dbs.es.nodeSelector }}
      - {{$key}}: {{$value}}
    {{- end }}
    hugePages:
      enabled: {{ .Values.dbs.pg.hugePages.enabled }}
      size: {{ .Values.dbs.pg.hugePages.size }}
      memory: "{{ .Values.dbs.pg.hugePages.memory }}"
  {{- if eq .Values.spec "ccp"  }}
  redis:
    enabled: {{ .Values.dbs.redis.enabled }}
    storageSize: {{ .Values.dbs.redis.storageSize }}
    storageClass: "{{ .Values.dbs.redis.storageClass }}"
    nodeSelector:
    {{- range $key, $value := .Values.dbs.redis.nodeSelector }}
      - {{$key}}: {{$value}}
    {{- end }}
  {{- end }}

{{- end }}
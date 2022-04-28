{{- define "spec.app_dbs" }}
dbs:
  cvat:
    enabled: {{ .Values.dbs.cvat.enabled }}
    pg:
      enabled: {{ .Values.dbs.cvat.enabled }}
    redis:
      enabled: {{ .Values.dbs.cvat.enabled }}
  es:
    enabled: {{ .Values.dbs.es.enabled }}
    storageSize: {{ .Values.dbs.es.storageSize }}
    storageClass: "{{ .Values.dbs.es.storageClass }}"
    patchEsNodes: {{ .Values.dbs.es.patchEsNodes }}
    {{- if .Values.dbs.es.nodeSelector }}
    nodeSelector:
    {{- range $key, $value := .Values.dbs.es.nodeSelector }}
      {{$key}}: {{$value}}
    {{- end }}
    {{- end }}
    cleanupPolicy:
        all: {{.Values.dbs.es.cleanupPolicy.all}}
        app: {{.Values.dbs.es.cleanupPolicy.app}}
        jobs: {{.Values.dbs.es.cleanupPolicy.jobs}}
        endpoints: {{.Values.dbs.es.cleanupPolicy.endpoints}}
  {{- if and (eq .Values.controlPlane.objectStorage.endpoint "") (eq .Values.controlPlane.objectStorage.type "minio")}}
  minio:
    enabled: {{ .Values.dbs.minio.enabled }}
    storageSize: {{ .Values.dbs.minio.storageSize }}
    storageClass: "{{ .Values.dbs.minio.storageClass }}"
    {{- if .Values.dbs.minio.nodeSelector }}
    nodeSelector:
    {{- range $key, $value := .Values.dbs.minio.nodeSelector }}
      {{$key}}: {{$value}}
    {{- end }}
    {{- end }}
  {{- end }}
  pg:
    enabled: {{ .Values.dbs.pg.enabled }}
    storageSize: {{ .Values.dbs.pg.storageSize }}
    storageClass: "{{ .Values.dbs.pg.storageClass }}"
    {{- if .Values.dbs.pg.nodeSelector }}
    nodeSelector:
    {{- range $key, $value := .Values.dbs.pg.nodeSelector }}
      {{$key}}: {{$value}}
    {{- end }}
    {{- end }}
    hugePages:
      enabled: {{ .Values.dbs.pg.hugePages.enabled }}
      size: {{ .Values.dbs.pg.hugePages.size }}
      memory: "{{ .Values.dbs.pg.hugePages.memory }}"
    backup:
      enabled: {{.Values.backup.enabled}}
      rotation: {{.Values.backup.rotation}}
      period: {{.Values.backup.period}}
  {{- if eq .Values.spec "ccp"  }}
  redis:
    enabled: {{ .Values.dbs.redis.enabled }}
    storageSize: {{ .Values.dbs.redis.storageSize }}
    storageClass: "{{ .Values.dbs.redis.storageClass }}"
    {{- if .Values.dbs.redis.nodeSelector }}
    nodeSelector:
    {{- range $key, $value := .Values.dbs.redis.nodeSelector }}
      {{$key}}: {{$value}}
    {{- end }}
    {{- end }}
  {{- end }}

{{- end }}
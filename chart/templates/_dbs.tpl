{{- define "spec.dbs" }}
dbs:
  es:
    cpuLimit: {{ .Values.dbs.es.cpuLimit }}
    cpuRequest: {{ .Values.dbs.es.cpuRequest }}
    enabled: {{ .Values.dbs.es.enabled }}
    fsGroup: {{ .Values.dbs.es.fsGroup }}
    image: {{ .Values.dbs.es.image }}
    memoryLimit: {{ .Values.dbs.es.memoryLimit }}
    memoryRequest: {{ .Values.dbs.es.memoryRequest }}
    nodePort: {{ .Values.dbs.es.nodePort }}
    patchEsNodes: {{ .Values.dbs.es.patchEsNodes }}
    port: {{ .Values.dbs.es.port }}
    runAsUser: {{ .Values.dbs.es.runAsUser }}
    serviceAccount: {{ .Values.dbs.es.serviceAccount }}
    storageSize: {{ .Values.dbs.es.storageSize }}
    svcName: {{ .Values.dbs.es.svcName }}
  minio:
    cpuRequest: {{ .Values.dbs.minio.cpuRequest }}
    enabled: {{ .Values.dbs.minio.enabled }}
    image: {{ .Values.dbs.minio.image }}
    memoryRequest: {{ .Values.dbs.minio.memoryRequest }}
    nodePort: {{ .Values.dbs.minio.nodePort }}
    port: {{ .Values.dbs.minio.port }}
    replicas: {{ .Values.dbs.minio.replicas }}
    serviceAccount: {{ .Values.dbs.minio.serviceAccount }}
    sharedStorage:
      consistentHash:
        key: {{ .Values.dbs.minio.sharedStorage.consistentHash.key }}
        value: {{ .Values.dbs.minio.sharedStorage.consistentHash.value }}
      enabled: {{ .Values.dbs.minio.sharedStorage.enabled }}
    storageSize: {{ .Values.dbs.minio.storageSize }}
    svcName: {{ .Values.dbs.minio.svcName }}
  pg:
    cpuRequest: {{ .Values.dbs.pg.cpuRequest }}
    dbname: {{ .Values.dbs.pg.dbname }}
    enabled: {{ .Values.dbs.pg.enabled }}
    fixpg: {{ .Values.dbs.pg.fixpg }}
    fsGroup: {{ .Values.dbs.pg.fsGroup }}
    hugePages:
      enabled: {{ .Values.dbs.pg.hugePages.enabled }}
      size: {{ .Values.dbs.pg.hugePages.size }}
      memory: {{ .Values.dbs.pg.hugePages.memory }}
    image: {{ .Values.dbs.pg.image }}
    maxConnections: {{ .Values.dbs.pg.maxConnections }}
    memoryRequest: {{ .Values.dbs.pg.memoryRequest }}
    pass: {{ .Values.dbs.pg.pass }}
    port: {{ .Values.dbs.pg.port }}
    runAsUser: {{ .Values.dbs.pg.runAsUser }}
    secretName: {{ .Values.dbs.pg.secretName }}
    serviceAccount: {{ .Values.dbs.pg.serviceAccount }}
    sharedBuffers: {{ .Values.dbs.pg.sharedBuffers }}
    storageSize: {{ .Values.dbs.pg.storageSize }}
    svcName: {{ .Values.dbs.pg.svcName }}
    user: {{ .Values.dbs.pg.user }}
  redis:
    appendonly: {{ .Values.dbs.redis.appendonly }}
    enabled: {{ .Values.dbs.redis.enabled }}
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
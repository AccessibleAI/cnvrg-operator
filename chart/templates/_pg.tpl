{{- define "spec.pg" }}
pg:
  enabled: "{{ .Values.pg.enabled }}"
  svcName: "{{ .Values.pg.svcName }}"
  image: "{{ .Values.pg.image }}"
  port: "{{ .Values.pg.port }}"
  dbname: "{{ .Values.pg.dbname }}"
  pass: "{{ .Values.pg.pass }}"
  user: "{{ .Values.pg.user }}"
  runAsUser: "{{ .Values.pg.runAsUser }}"
  runAsGroup: "{{ .Values.pg.runAsGroup }}"
  fsGroup: "{{ .Values.pg.fsGroup }}"
  storageClass: "{{ .Values.pg.storageClass }}"
  hugePages:
    enabled: "{{.Values.pg.hugePages.enabled}}"
    size: "{{.Values.pg.hugePages.size}}"
    memory: "{{.Values.pg.hugePages.memory}}"

  {{- if eq .Values.computeProfile "large"}}
  cpuRequest: "{{ .Values.computeProfiles.large.pg.cpu }}"
  memoryRequest: "{{ .Values.computeProfiles.large.pg.memory }}"
  storageSize: "{{ .Values.computeProfiles.large.storage }}"
  {{- end }}

  {{- if eq .Values.computeProfile "medium"}}
  cpuRequest: "{{ .Values.computeProfiles.medium.pg.cpu }}"
  memoryRequest: "{{ .Values.computeProfiles.medium.pg.memory }}"
  storageSize: "{{ .Values.computeProfiles.medium.storage }}"
  {{- end }}

  {{- if eq .Values.computeProfile "small"}}
  cpuRequest: "{{ .Values.computeProfiles.small.pg.cpu }}"
  memoryRequest: "{{ .Values.computeProfiles.small.pg.memory }}"
  storageSize: "{{ .Values.computeProfiles.small.storage }}"
  {{- end }}

pgBackup:
  enabled: "{{ .Values.pgBackup.enabled }}"
  name: "{{ .Values.pgBackup.name }}"
  path: "{{ .Values.pgBackup.path }}"
  scriptPath: "{{ .Values.pgBackup.scriptPath }}"
  storageClass: "{{ .Values.pgBackup.storageClass }}"
  cronTime: "{{ .Values.pgBackup.cronTime }}"

  {{- if eq .Values.computeProfile "large"}}
  storageSize: "{{ .Values.computeProfiles.large.storage }}"
  {{- end }}

  {{- if eq .Values.computeProfile "medium"}}
  storageSize: "{{ .Values.computeProfiles.medium.storage }}"
  {{- end }}

  {{- if eq .Values.computeProfile "small"}}
  storageSize: "{{ .Values.computeProfiles.small.storage }}"
  {{- end }}

{{- end }}
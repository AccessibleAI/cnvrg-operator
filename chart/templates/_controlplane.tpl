{{- define "spec.controlPlane" }}
controlPlane:
  baseConfig:
    agentCustomTag: {{ .Values.controlPlane.baseConfig.agentCustomTag }}
    intercom: "{{ .Values.controlPlane.baseConfig.intercom }}"
    featureFlags:
    {{- range $key, $value := .Values.controlPlane.baseConfig.featureFlags }}
      {{$key}}: {{$value}}
    {{- end }}
  hyper:
    enabled: {{ .Values.controlPlane.hyper.enabled }}
  objectStorage:
    cnvrgStorageAccessKey: {{ .Values.controlPlane.objectStorage.cnvrgStorageAccessKey }}
    cnvrgStorageBucket: {{ .Values.controlPlane.objectStorage.cnvrgStorageBucket }}
    cnvrgStorageRegion: {{ .Values.controlPlane.objectStorage.cnvrgStorageRegion }}
    cnvrgStorageSecretKey: {{ .Values.controlPlane.objectStorage.cnvrgStorageSecretKey }}
    cnvrgStorageType: {{ .Values.controlPlane.objectStorage.cnvrgStorageType }}
    gcpKeyfileMountPath: {{ .Values.controlPlane.objectStorage.gcpKeyfileMountPath }}
    gcpKeyfileName: {{ .Values.controlPlane.objectStorage.gcpKeyfileName }}
    gcpStorageSecret: {{ .Values.controlPlane.objectStorage.gcpStorageSecret }}
  searchkiq:
    enabled: {{ .Values.controlPlane.searchkiq.enabled }}
  sidekiq:
    enabled: {{ .Values.controlPlane.sidekiq.enabled }}
    split: {{ .Values.controlPlane.sidekiq.split }}
  smtp:
    server: "{{ .Values.controlPlane.smtp.server }}"
    port: {{ .Values.controlPlane.smtp.port }}
    username: "{{ .Values.controlPlane.smtp.username }}"
    password: "{{ .Values.controlPlane.smtp.password }}"
    domain: "{{ .Values.controlPlane.smtp.domain }}"
  systemkiq:
    enabled: {{ .Values.controlPlane.systemkiq.enabled }}
  webapp:
    enabled: {{ .Values.controlPlane.webapp.enabled }}
    image: {{ .Values.controlPlane.webapp.image }}
    replicas: {{ .Values.controlPlane.webapp.replicas }}
  mpi:
    enabled: {{ .Values.controlPlane.mpi.enabled }}
    image: {{ .Values.controlPlane.mpi.image }}
    kubectlDeliveryImage: {{ .Values.controlPlane.mpi.kubectlDeliveryImage }}
    extraArgs:
    {{- range $key, $value := .Values.controlPlane.mpi.extraArgs }}
      {{$key}}: {{$value}}
    {{- end }}
{{- end }}
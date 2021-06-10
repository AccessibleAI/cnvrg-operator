{{- define "spec.controlPlane" }}
controlPlane:
  image: {{ .Values.controlPlane.image }}
  baseConfig:
    agentCustomTag: {{ .Values.controlPlane.baseConfig.agentCustomTag }}
    intercom: "{{ .Values.controlPlane.baseConfig.intercom }}"
    {{- if .Values.controlPlane.baseConfig.featureFlags }}
    featureFlags:
    {{- range $key, $value := .Values.controlPlane.baseConfig.featureFlags }}
      {{$key}}: "{{$value}}"
    {{- end }}
    {{- else }}
    featureFlags: { }
    {{- end }}
  hyper:
    enabled: {{ .Values.controlPlane.hyper.enabled }}
  cnvrgRouter:
    enabled: {{ .Values.controlPlane.cnvrgRouter.enabled }}
    image: {{ .Values.controlPlane.cnvrgRouter.image }}
  objectStorage:
    accessKey: "{{ .Values.controlPlane.objectStorage.accessKey }}"
    bucket: "{{ .Values.controlPlane.objectStorage.bucket }}"
    region: "{{ .Values.controlPlane.objectStorage.region }}"
    secretKey: "{{ .Values.controlPlane.objectStorage.secretKey }}"
    type: {{ .Values.controlPlane.objectStorage.type }}
    endpoint: "{{ .Values.controlPlane.objectStorage.endpoint }}"
    azureAccountName: "{{ .Values.controlPlane.objectStorage.azureAccountName}}"
    azureContainer: "{{.Values.controlPlane.objectStorage.azureContainer}}"
    gcpSecretRef: "{{ .Values.controlPlane.objectStorage.gcpSecretRef }}"
    gcpProject: "{{ .Values.controlPlane.objectStorage.gcpProject }}"
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
    replicas: {{ .Values.controlPlane.webapp.replicas }}
  mpi:
    enabled: {{ .Values.controlPlane.mpi.enabled }}
    image: {{ .Values.controlPlane.mpi.image }}
    kubectlDeliveryImage: {{ .Values.controlPlane.mpi.kubectlDeliveryImage }}
    {{- if .Values.controlPlane.mpi.extraArgs }}
    extraArgs:
    {{- range $key, $value := .Values.controlPlane.mpi.extraArgs }}
      {{$key}}: {{$value}}
    {{- end }}
    {{- else }}
    extraArgs: { }
    {{- end }}
{{- end }}
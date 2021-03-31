{{- define "spec.minio" }}
minio:
  enabled: "{{ .Values.minio.enabled }}"
  replicas: "{{ .Values.minio.replicas }}"
  image: "{{ .Values.minio.image }}"
  port: "{{ .Values.minio.port }}"
  svcName: "{{ .Values.minio.svcName }}"
  nodePort: "{{ .Values.minio.nodePort }}"
  storageClass: "{{ .Values.minio.storageClass }}"
  sharedStorage:
    enabled: "{{ .Values.minio.sharedStorage.enabled }}"


  {{- if eq .Values.computeProfile "large"}}
  cpuRequest: "{{ .Values.computeProfiles.large.minio.cpu }}"
  memoryRequest: "{{ .Values.computeProfiles.large.minio.memory }}"
  storageSize: "{{ .Values.computeProfiles.large.storage }}"
  {{- end }}

  {{- if eq .Values.computeProfile "medium"}}
  cpuRequest: "{{ .Values.computeProfiles.medium.minio.cpu }}"
  memoryRequest: "{{ .Values.computeProfiles.medium.minio.memory }}"
  storageSize: "{{ .Values.computeProfiles.medium.storage }}"
  {{- end }}

  {{- if eq .Values.computeProfile "small"}}
  cpuRequest: "{{ .Values.computeProfiles.small.minio.cpu }}"
  memoryRequest: "{{ .Values.computeProfiles.small.minio.memory }}"
  storageSize: "{{ .Values.computeProfiles.small.storage }}"
  {{- end }}

{{- end }}
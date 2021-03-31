{{- define "spec.logging" }}
logging:
  enabled: "{{ .Values.logging.enabled }}"
  es:
    enabled: "{{ .Values.logging.es.enabled }}"
    image: "{{ .Values.logging.es.image }}"
    maxMapImage: "{{.Values.logging.es.maxMapImage}}"
    port: "{{ .Values.logging.es.port }}"
    svcName: "{{ .Values.logging.es.svcName }}"
    runAsUser: "{{ .Values.logging.es.runAsUser }}"
    runAsGroup: "{{ .Values.logging.es.runAsGroup }}"
    fsGroup: "{{ .Values.logging.es.fsGroup }}"
    patchEsNodes: "{{ .Values.logging.es.patchEsNodes }}"
    nodePort: "{{ .Values.logging.es.nodePort }}"
    storageClass: "{{ .Values.logging.es.storageClass }}"
    javaOpts: "{{ .Values.logging.es.javaOpts }}"
    cpuLimit: "{{.Values.logging.es.cpuLimit}}"
    memoryLimit: "{{.Values.logging.es.memoryLimit}}"
    {{- if eq .Values.computeProfile "large"}}
    cpuRequest: "{{ .Values.computeProfiles.large.es.cpu }}"
    memoryRequest: "{{ .Values.computeProfiles.large.es.memory }}"
    storageSize: "{{ .Values.computeProfiles.large.storage }}"
    {{- end }}
    {{- if eq .Values.computeProfile "medium"}}
    cpuRequest: "{{ .Values.computeProfiles.medium.es.cpu }}"
    memoryRequest: "{{ .Values.computeProfiles.medium.es.memory }}"
    storageSize: "{{ .Values.computeProfiles.medium.storage }}"
    {{- end }}
    {{- if eq .Values.computeProfile "small"}}
    cpuRequest: "{{ .Values.computeProfiles.small.es.cpu }}"
    memoryRequest: "{{ .Values.computeProfiles.small.es.memory }}"
    storageSize: "{{ .Values.computeProfiles.small.storage }}"
    {{- end }}
  elastalert:
    enabled: "{{ .Values.logging.elastalert.enabled }}"
    image: "{{ .Values.logging.elastalert.image }}"
    port: "{{ .Values.logging.elastalert.port }}"
    nodePort: "{{ .Values.logging.elastalert.nodePort }}"
    containerPort: "{{ .Values.logging.elastalert.containerPort }}"
    svcName: "{{ .Values.logging.elastalert.svcName }}"
    storageClass: "{{ .Values.logging.elastalert.storageClass }}"
    runAsUser: "{{ .Values.logging.elastalert.runAsUser }}"
    runAsGroup: "{{ .Values.logging.elastalert.runAsGroup }}"
    fsGroup: "{{ .Values.logging.elastalert.fsGroup }}"
    {{- if eq .Values.computeProfile "large"}}
    cpuRequest: "{{ .Values.computeProfiles.large.elastalert.cpu }}"
    memoryRequest: "{{ .Values.computeProfiles.large.elastalert.memory }}"
    storageSize: "{{ .Values.computeProfiles.large.storage }}"
    {{- end }}
    {{- if eq .Values.computeProfile "medium"}}
    cpuRequest: "{{ .Values.computeProfiles.medium.elastalert.cpu }}"
    memoryRequest: "{{ .Values.computeProfiles.medium.elastalert.memory }}"
    storageSize: "{{ .Values.computeProfiles.medium.storage }}"
    {{- end }}
    {{- if eq .Values.computeProfile "small"}}
    cpuRequest: "{{ .Values.computeProfiles.small.elastalert.cpu }}"
    memoryRequest: "{{ .Values.computeProfiles.small.elastalert.memory }}"
    storageSize: "{{ .Values.computeProfiles.small.storage }}"
    cpuLimit: "{{ .Values.computeProfiles.small.elastalert.cpu }}"
    memoryLimit: "{{ .Values.computeProfiles.small.elastalert.memory }}"
    {{- end }}
  fluentd:
    enabled: "{{ .Values.logging.fluentd.enabled }}"
    image: "{{.Values.logging.fluentd.image}}"
    journalPath: "{{ .Values.logging.fluentd.journalPath }}"
    containersPath: "{{ .Values.logging.fluentd.containersPath }}"
    journald: "{{ .Values.logging.fluentd.journald }}"
    cpuRequest: "{{ .Values.logging.fluentd.cpuRequest }}"
    memoryRequest: "{{ .Values.logging.fluentd.memoryRequest }}"
    memoryLimit: "{{ .Values.logging.fluentd.memoryLimit }}"
  kibana:
    enabled: "{{ .Values.logging.kibana.enabled }}"
    svcName: "{{ .Values.logging.kibana.svcName }}"
    image: "{{ .Values.logging.kibana.image }}"
    nodePort: "{{ .Values.logging.kibana.nodePort }}"
    cpuRequest: "{{ .Values.logging.kibana.cpuRequest }}"
    memoryRequest: "{{ .Values.logging.kibana.memoryRequest }}"
    cpuLimit: "{{ .Values.logging.kibana.cpuLimit }}"
    memoryLimit: "{{ .Values.logging.kibana.memoryLimit }}"
    port: "{{ .Values.logging.kibana.port }}"
{{- end }}
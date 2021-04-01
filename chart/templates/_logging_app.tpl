{{- define "spec.logging_app" }}

logging:
  enabled: "{{ .Values.logging.enabled }}"
  elastalert:
    enabled: "{{ .Values.logging.elastalert.enabled }}"
    image: {{ .Values.logging.elastalert.image }}
    port: {{ .Values.logging.elastalert.port }}
    nodePort: {{ .Values.logging.elastalert.nodePort }}
    containerPort: {{ .Values.logging.elastalert.containerPort }}
    storageSize: {{ .Values.logging.elastalert.storageSize }}
    svcName: {{ .Values.logging.elastalert.svcName }}
    storageClass: "{{ .Values.logging.elastalert.storageClass }}"
    cpuRequest: {{ .Values.logging.elastalert.cpuRequest }}
    memoryRequest: {{ .Values.logging.elastalert.memoryRequest }}
    cpuLimit: {{ .Values.logging.elastalert.cpuLimit }}
    memoryLimit: {{ .Values.logging.elastalert.memoryLimit }}
    runAsUser: {{ .Values.logging.elastalert.runAsUser }}
    fsGroup: {{ .Values.logging.elastalert.fsGroup }}
  kibana:
    enabled: "{{ .Values.logging.kibana.enabled }}"
    serviceAccount: {{ .Values.logging.kibana.serviceAccount }}
    svcName: {{ .Values.logging.kibana.svcName }}
    port: {{ .Values.logging.kibana.port }}
    image: {{ .Values.logging.kibana.image }}
    nodePort: {{ .Values.logging.kibana.nodePort }}
    cpuRequest: {{ .Values.logging.kibana.cpuRequest }}
    memoryRequest: {{ .Values.logging.kibana.memoryRequest }}
    cpuLimit: {{ .Values.logging.kibana.cpuLimit }}
    memoryLimit: {{ .Values.logging.kibana.memoryLimit }}
{{- end }}
{{- define "spec.logging_app" }}
logging:
  elastalert:
    enabled: {{ .Values.logging.elastalert.enabled }}
    storageSize: {{ .Values.logging.elastalert.storageSize }}
    storageClass: "{{ .Values.logging.elastalert.storageClass }}"
    {{- if .Values.logging.elastalert.nodeSelector }}
    nodeSelector:
    {{- range $key, $value := .Values.logging.elastalert.nodeSelector }}
      {{$key}}: {{$value}}
    {{- end }}
    {{- end }}
  kibana:
    enabled: {{ .Values.logging.kibana.enabled }}
    svcName: {{ .Values.logging.kibana.svcName }}
{{- end }}
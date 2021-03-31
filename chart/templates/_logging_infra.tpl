{{- define "spec.logging_infra" }}
logging:
  enabled: {{ .Values.logging.enabled }}
  fluentbit:
    image: {{ .Values.logging.fluentbit.image }}
{{- end }}
{{- define "spec.logging_infra" }}
logging:
  fluentbit:
    enabled: {{ .Values.logging.fluentbit.enabled }}
{{- end }}
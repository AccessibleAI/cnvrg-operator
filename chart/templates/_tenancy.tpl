{{- define "spec.tenancy" }}
tenancy:
  enabled: {{ .Values.tenancy.enabled }}
  key: {{ .Values.tenancy.key }}
  value: {{ .Values.tenancy.value }}
{{- end }}
{{- define "spec.labelsAndAnnotations" }}

{{- if .Values.labels }}
labels:
  {{- range $key, $value := .Values.labels }}
  {{ $key }}: {{ $value | quote }}
  {{- end }}
{{- else }}
labels: { }
{{- end }}

{{- if .Values.annotations }}
annotations:
  {{- range $key, $value := .Values.annotations }}
  {{ $key }}: {{ $value | quote }}
  {{- end }}
{{- else }}
annotations: { }
{{- end }}

{{- end }}

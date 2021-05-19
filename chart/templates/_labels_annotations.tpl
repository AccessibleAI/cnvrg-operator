{{- define "spec.labelsAndAnnotations" }}

{{- if .Values.labels }}
labels:
  {{- range $key, $value := .Values.labels }}
  {{$key}}: {{$value}}
  {{- end }}
{{- else }}
labels: { }
{{- end }}

{{- if .Values.annotations}}
annotations:
  {{- range $key, $value := .Values.annotations }}
  {{$key}}: {{$value}}
  {{- end }}
{{- else }}
annotations: { }
{{- end }}

{{- end }}
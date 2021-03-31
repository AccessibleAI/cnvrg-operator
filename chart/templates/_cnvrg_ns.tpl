{{- define "spec.cnvrgNs" -}}
{{- if eq .Release.Namespace "default"}}cnvrg{{- else -}}{{ .Release.Namespace }}{{- end -}}
{{- end -}}

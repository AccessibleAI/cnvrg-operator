{{- define "spec.redis" }}
redis:
  enabled: "{{ .Values.redis.enabled }}"
  image: "{{ .Values.redis.image }}"
  svcName: "{{ .Values.redis.svcName }}"
  port: "{{ .Values.redis.port }}"
  limits:
    cpu: "{{.Values.redis.limits.cpu}}"
    memory: "{{.Values.redis.limits.memory}}"

  {{- if eq .Values.computeProfile "large" }}
  requests:
    cpu: "{{ .Values.computeProfiles.large.redis.cpu }}"
    memory:  "{{ .Values.computeProfiles.large.redis.memory }}"
  {{- end }}

  {{- if eq .Values.computeProfile "medium" }}
  requests:
    cpu: "{{ .Values.computeProfiles.medium.redis.cpu }}"
    memory:  "{{ .Values.computeProfiles.medium.redis.memory }}"
  {{- end }}

  {{- if eq .Values.computeProfile "small" }}
  requests:
    cpu: "{{ .Values.computeProfiles.small.redis.cpu }}"
    memory:  "{{ .Values.computeProfiles.small.redis.memory }}"
  {{- end }}

{{- end }}
{{- define "spec.registry" }}
registry:
  url: "{{.Values.registry.url}}"
  user: "{{.Values.registry.user}}"
  password: "{{.Values.registry.password}}"
{{- end }}
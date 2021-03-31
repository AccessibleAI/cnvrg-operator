{{- define "spec.registry" }}
registry:
  name: {{.Values.registry.name}}
  url: {{.Values.registry.url}}
  user: {{.Values.registry.user}}
  password: {{.Values.registry.password}}
{{- end }}
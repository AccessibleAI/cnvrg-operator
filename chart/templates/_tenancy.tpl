{{- define "spec.tenancy" }}
tenancy:
  enabled: "{{.Values.tenancy.enabled}}"
  dedicatedNodes: "{{.Values.tenancy.dedicatedNodes}}"
  cnvrg:
    key: "{{.Values.tenancy.cnvrg.key}}"
    value: "{{.Values.tenancy.cnvrg.value}}"
{{- end }}

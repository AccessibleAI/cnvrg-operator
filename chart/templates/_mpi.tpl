{{- define "spec.mpi" }}
mpi:
  enabled: "{{ .Values.mpi.enabled }}"
  image: "{{.Values.mpi.image}}"
  kubectlDeliveryImage: "{{.Values.mpi.kubectlDeliveryImage}}"
  registry:
    name: "{{.Values.mpi.registry.name}}"
    url: "{{.Values.mpi.registry.url}}"
    user: "{{.Values.mpi.registry.user}}"
    password: "{{.Values.mpi.registry.password}}"
{{- end }}
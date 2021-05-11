{{- define "spec.storage"  }}
storage:
  hostpath:
    enabled: {{ .Values.storage.hostpath.enabled }}
    image: "{{.Values.storage.hostpath.image}}"
    path: "{{ .Values.storage.hostpath.path }}"
    reclaimPolicy: "{{.Values.storage.hostpath.reclaimPolicy}}"
    defaultSc: {{.Values.storage.hostpath.defaultSc}}
    {{- if .Values.storage.hostpath.nodeSelector }}
    nodeSelector:
    {{- range $key, $value := .Values.storage.hostpath.nodeSelector }}
      {{$key}}: {{$value}}
    {{- end }}
    {{- end }}
  nfs:
    enabled: {{ .Values.storage.nfs.enabled }}
    image: "{{.Values.storage.nfs.image}}"
    server: "{{ .Values.storage.nfs.server }}"
    path: "{{ .Values.storage.nfs.path }}"
    reclaimPolicy: "{{.Values.storage.nfs.reclaimPolicy}}"
    defaultSc: {{.Values.storage.nfs.defaultSc}}
{{- end }}
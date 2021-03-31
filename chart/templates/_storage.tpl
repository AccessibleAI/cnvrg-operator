{{- define "spec.storage"  }}
storage:
  enabled: "{{.Values.storage.enabled}}"
  hostpath:
    enabled: "{{ .Values.storage.hostpath.enabled }}"
    image: "{{.Values.storage.hostpath.image}}"
    hostPath: "{{ .Values.storage.hostpath.hostPath }}"
    storageClassName: "{{ .Values.storage.hostpath.storageClassName }}"
    nodeName: "{{ .Values.storage.hostpath.nodeName }}"
    cpuRequest: "{{ .Values.storage.hostpath.cpuRequest }}"
    memoryRequest: "{{ .Values.storage.hostpath.memoryRequest }}"
    cpuLimit: "{{ .Values.storage.hostpath.cpuLimit }}"
    memoryLimit: "{{ .Values.storage.hostpath.memoryLimit }}"
    reclaimPolicy: "{{.Values.storage.hostpath.reclaimPolicy}}"
    defaultSc: "{{.Values.storage.hostpath.defaultSc}}"
  nfs:
    enabled: "{{ .Values.storage.nfs.enabled }}"
    image: "{{.Values.storage.nfs.image}}"
    provisioner: "{{ .Values.storage.nfs.provisioner }}"
    storageClassName: "{{ .Values.storage.nfs.storageClassName }}"
    server: "{{ .Values.storage.nfs.server }}"
    path: "{{ .Values.storage.nfs.path }}"
    cpuRequest: "{{ .Values.storage.nfs.cpuRequest }}"
    memoryRequest: "{{ .Values.storage.nfs.memoryRequest }}"
    cpuLimit: "{{ .Values.storage.nfs.cpuLimit }}"
    memoryLimit: "{{ .Values.storage.nfs.memoryLimit }}"
    reclaimPolicy: "{{.Values.storage.nfs.reclaimPolicy}}"
    defaultSc: "{{.Values.storage.nfs.defaultSc}}"
{{- end }}
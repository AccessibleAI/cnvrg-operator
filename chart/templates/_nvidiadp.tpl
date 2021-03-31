{{- define "spec.nvidiadp" }}
nvidiadp:
  enabled: "{{ .Values.nvidiadp.enabled }}"
  image: "{{ .Values.nvidiadp.image }}"
  nodeSelector:
    enabled: "{{ .Values.nvidiadp.nodeSelector.enabled }}"
    key: "{{ .Values.nvidiadp.nodeSelector.key }}"
    value: "{{ .Values.nvidiadp.nodeSelector.value }}"
{{- end }}

{{- define "spec.globals"  }}
debug: "{{ .Values.debug }}"
otags: "{{ .Values.otags }}"
dumpDir: "{{ .Values.dumpDir }}"
dryRun: "{{ .Values.dryRun }}"
clusterDomain: "{{ .Values.clusterDomain }}"
{{- end }}

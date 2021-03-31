{{- define "spec.vpa" }}
vpa:
  enabled: "false"
  images:
    admissionImage: "{{.Values.vpa.images.admissionImage}}"
    recommenderImage: "{{.Values.vpa.images.recommenderImage}}"
    updaterImage: "{{.Values.vpa.images.updaterImage}}"
{{- end }}

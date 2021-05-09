{{- define "spec.networking_app" }}
networking:
  https:
    enabled: {{ .Values.networking.https.enabled }}
    cert: "{{ .Values.networking.https.cert }}"
    key: "{{ .Values.networking.https.key }}"
    certSecret: "{{ .Values.networking.https.certSecret }}"
  ingress:
    type: {{ .Values.networking.ingress.type }}
{{- end }}

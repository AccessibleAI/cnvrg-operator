{{- define "spec.networking_app" }}
networking:
  https:
    enabled: "{{ .Values.networking.https.enabled }}"
    cert: {{ .Values.networking.https.cert }}
    key: {{ .Values.networking.https.key }}
    certSecret: {{ .Values.networking.https.certSecret }}
  ingress:
    ingressType: {{ .Values.networking.ingress.ingressType }}
    perTryTimeout: {{ .Values.networking.ingress.perTryTimeout }}
    retriesAttempts: {{ .Values.networking.ingress.retriesAttempts }}
    timeout: {{ .Values.networking.ingress.timeout }}
{{- end }}

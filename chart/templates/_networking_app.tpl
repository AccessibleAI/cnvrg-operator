{{- define "spec.networking_app" }}
networking:
  https:
    enabled: {{ .Values.networking.https.enabled }}
    certSecret: "{{ .Values.networking.https.certSecret }}"
  ingress:
    type: {{ .Values.networking.ingress.type }}
    istioGwName: "{{.Values.networking.ingress.istioGwName}}"
    {{- if and (eq .Values.networking.ingress.type "istio") (eq .Values.spec "allinone") }}
    istioGwEnabled: false
    {{- else }}
    istioGwEnabled: {{.Values.networking.ingress.istioGwEnabled}}
    {{- end }}
{{- end }}

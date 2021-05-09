{{- define "spec.networking_infra" }}
networking:
  https:
    enabled: {{ .Values.networking.https.enabled }}
    cert: "{{ .Values.networking.https.cert }}"
    key: "{{ .Values.networking.https.key }}"
    certSecret: "{{ .Values.networking.https.certSecret }}"
  ingress:
    type: "{{ .Values.networking.ingress.type }}"
  istio:
    enabled: {{ .Values.networking.istio.enabled }}
    externalIp:
    {{- range $_, $value := .Values.networking.istio.externalIp }}
      - {{$value}}
    {{- end }}
    ingressSvcAnnotations:
    {{- range $key, $value := .Values.networking.istio.ingressSvcAnnotations }}
      {{$key}}: {{$value}}
    {{- end }}
    ingressSvcExtraPorts:
    {{- range $_, $value := .Values.networking.istio.ingressSvcExtraPorts }}
      - {{$value}}
    {{- end }}
    lbSourceRanges:
    {{- range $_, $value := .Values.networking.istio.lbSourceRanges }}
      - {{$value}}
    {{- end }}
{{- end }}

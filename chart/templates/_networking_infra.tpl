{{- define "spec.networking_infra" }}
networking:
  https:
    enabled: {{ .Values.networking.https.enabled }}
    cert: "{{ .Values.networking.https.cert }}"
    key: "{{ .Values.networking.https.key }}"
    certSecret: "{{ .Values.networking.https.certSecret }}"
  ingress:
    type: "{{ .Values.networking.ingress.type }}"
    istioGwEnabled: {{.Values.networking.ingress.istioGwEnabled}}
  istio:
    enabled: {{ .Values.networking.istio.enabled }}
    {{- if .Values.networking.istio.externalIp }}
    externalIp:
    {{- range $_, $value := .Values.networking.istio.externalIp }}
      - {{$value}}
    {{- end }}
    {{- end }}
    {{- if .Values.networking.istio.ingressSvcAnnotations }}
    ingressSvcAnnotations:
    {{- range $key, $value := .Values.networking.istio.ingressSvcAnnotations }}
      {{$key}}: {{$value}}
    {{- end }}
    {{- end }}
    {{- if .Values.networking.istio.ingressSvcExtraPorts }}
    ingressSvcExtraPorts:
    {{- range $_, $value := .Values.networking.istio.ingressSvcExtraPorts }}
      - {{$value}}
    {{- end }}
    {{- end }}
    {{- if .Values.networking.istio.lbSourceRanges }}
    lbSourceRanges:
    {{- range $_, $value := .Values.networking.istio.lbSourceRanges }}
      - {{$value}}
    {{- end }}
    {{- end }}
{{- end }}

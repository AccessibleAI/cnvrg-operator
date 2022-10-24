{{- define "spec.networking_infra" }}
networking:
  https:
    enabled: {{ .Values.networking.https.enabled }}
    certSecret: "{{ .Values.networking.https.certSecret }}"
  ingress:
    type: "{{ .Values.networking.ingress.type }}"
    {{- if eq .Values.networking.ingress.type "istio" }}
    istioGwEnabled: {{.Values.networking.ingress.istioGwEnabled}}
    istioGwName: "{{.Values.networking.ingress.istioGwName}}"
    {{- end }}
  proxy:
    enabled: {{ .Values.networking.proxy.enabled }}
    {{- if .Values.networking.proxy.httpProxy }}
    httpProxy:
    {{- range $_, $value := .Values.networking.proxy.httpProxy }}
      - {{$value}}
    {{- end }}
    {{- end }}
    {{- if .Values.networking.proxy.httpsProxy }}
    httpsProxy:
    {{- range $_, $value := .Values.networking.proxy.httpsProxy }}
      - {{$value}}
    {{- end }}
    {{- end }}
    {{- if .Values.networking.proxy.noProxy }}
    noProxy:
    {{- range $_, $value := .Values.networking.proxy.noProxy }}
      - {{$value}}
    {{- end }}
    {{- end }}
  {{- if eq .Values.networking.ingress.type "istio" }}
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
  eastWest:
    enabled: {{ .Values.networking.eastWest.enabled }}
    primary: {{ .Values.networking.eastWest.primary }}
    clusterName: {{ .Values.networking.eastWest.clusterName }}
    network: {{ .Values.networking.eastWest.network }}
    meshId: {{ .Values.networking.eastWest.meshId }}
{{- end }}

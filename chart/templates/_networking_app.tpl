{{- define "spec.networking_app" }}
networking:
  https:
    enabled: {{ .Values.networking.https.enabled }}
    certSecret: "{{ .Values.networking.https.certSecret }}"
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
  ingress:
    type: {{ .Values.networking.ingress.type }}
    {{- if eq .Values.networking.ingress.type "istio" }}
    istioGwName: "{{.Values.networking.ingress.istioGwName}}"
    {{- if and (eq .Values.networking.ingress.type "istio") (eq .Values.spec "allinone") }}
    istioGwEnabled: false
    {{- else }}
    istioGwEnabled: {{.Values.networking.ingress.istioGwEnabled}}
    {{- end }}
    {{- end }}
    eastWest:
      enabled: {{ .Values.networking.istio.eastWest.enabled }}
      primary: {{ .Values.networking.istio.eastWest.primary }}
      clusterName: {{ .Values.networking.istio.eastWest.clusterName }}
      network: {{ .Values.networking.istio.eastWest.network }}
      meshId: {{ .Values.networking.istio.eastWest.meshId }}
{{- end }}

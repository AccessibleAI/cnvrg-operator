{{- define "spec.networking_infra" }}
networking:
  https:
    enabled: "{{ .Values.networking.https.enabled }}"
    cert: "{{ .Values.networking.https.cert }}"
    key: "{{ .Values.networking.https.key }}"
    certSecret: "{{ .Values.networking.https.certSecret }}"
  ingress:
    ingressType: {{ .Values.networking.ingress.ingressType }}
    perTryTimeout: {{ .Values.networking.ingress.perTryTimeout }}
    retriesAttempts: {{ .Values.networking.ingress.retriesAttempts }}
    timeout: {{ .Values.networking.ingress.timeout }}
  istio:
    enabled: "{{ .Values.networking.istio.enabled }}"
    operatorImage: {{ .Values.networking.istio.operatorImage }}
    hub: {{ .Values.networking.istio.hub }}
    tag: {{ .Values.networking.istio.tag }}
    proxyImage: {{ .Values.networking.istio.proxyImage }}
    mixerImage: {{ .Values.networking.istio.mixerImage }}
    pilotImage: {{ .Values.networking.istio.pilotImage }}
    externalIp: "{{ .Values.networking.istio.externalIp }}"
    ingressSvcAnnotations: "{{ .Values.networking.istio.ingressSvcAnnotations }}"
    ingressSvcExtraPorts: "{{ .Values.networking.istio.ingressSvcExtraPorts }}"
    loadBalancerSourceRanges: "{{ .Values.networking.istio.loadBalancerSourceRanges }}"
{{- end }}

{{- define "spec.networking" }}
networking:
  enabled: "{{ .Values.networking.enabled}}"
  ingressType: "{{.Values.networking.ingressType}}"
  istio:
    enabled: "{{ .Values.networking.istio.enabled }}"
    operatorImage: "{{ .Values.networking.istio.operatorImage }}"
    hub: "{{.Values.networking.istio.hub}}"
    tag: "{{.Values.networking.istio.tag}}"
    proxyImage: "{{.Values.networking.istio.proxyImage}}"
    mixerImage: "{{.Values.networking.istio.mixerImage}}"
    pilotImage: "{{.Values.networking.istio.pilotImage}}"
    gwName: "{{ .Values.networking.istio.gwName }}"
    externalIp: "{{ .Values.networking.istio.externalIp }}"
    ingressSvcAnnotations: "{{.Values.networking.istio.ingressSvcAnnotations}}"
    loadBalancerSourceRanges: "{{.Values.networking.istio.loadBalancerSourceRanges}}"
  ingress:
    enabled: "{{ .Values.networking.ingress.enabled }}"
    timeout: "{{ .Values.networking.ingress.timeout }}"
    retriesAttempts: "{{ .Values.networking.ingress.retriesAttempts }}"
    perTryTimeout: "{{ .Values.networking.ingress.perTryTimeout }}"
  https:
    enabled: "{{ .Values.networking.https.enabled }}"
    cert: "{{ .Values.networking.https.cert }}"
    key: "{{ .Values.networking.https.key }}"
    certSecret: "{{ .Values.networking.https.certSecret }}"
{{- end }}

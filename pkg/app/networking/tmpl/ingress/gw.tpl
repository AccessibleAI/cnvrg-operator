apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: {{ .Spec.Networking.Ingress.IstioGwName }}
  namespace: {{ ns . }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  selector:
    istio: cnvrg-ingressgateway
  servers:
    - port:
        number: 80
        name: http
        protocol: HTTP
      hosts:
        - "*.{{ .Spec.ClusterDomain }}"
      {{- if and ( isTrue .Spec.Networking.HTTPS.Enabled ) (ne .Spec.Networking.HTTPS.CertSecret "") }}
      tls:
        httpsRedirect: true
    - hosts:
        - "*.{{ .Spec.ClusterDomain }}"
      port:
        name: https
        number: 443
        protocol: HTTPS
      tls:
        mode: SIMPLE
        credentialName: {{ .Spec.Networking.HTTPS.CertSecret }}
      {{- else if isTrue .Spec.Networking.HTTPS.Enabled }}
      tls:
        httpsRedirect: true
    - hosts:
        - "*.{{ .Spec.ClusterDomain }}"
      port:
        name: https
        number: 443
        protocol: HTTP
      {{- end }}

apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: {{ .Networking.Istio.GwName }}
  namespace: {{ .CnvrgNs }}
spec:
  selector:
    istio: ingressgateway
  servers:
    - port:
        number: 80
        name: http
        protocol: HTTP
      hosts:
        - "*.{{ .ClusterDomain }}"
      {{- if and (eq .Networking.HTTPS.Enabled "true") (ne .Networking.HTTPS.CertSecret "") }}
      tls:
        httpsRedirect: true
    - hosts:
        - "*.{{ .ClusterDomain }}"
      port:
        name: https
        number: 443
        protocol: HTTPS
      tls:
        mode: SIMPLE
        credentialName: "{{ .Networking.HTTPS.CertSecret }}"
      {{- else if eq .Networking.HTTPS.Enabled "true" }}
      tls:
        httpsRedirect: true
    - hosts:
        - "*.{{ .ClusterDomain }}"
      port:
        name: https
        number: 443
        protocol: HTTP
      {{- end }}

apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{.Logging.Es.SvcName }}
  namespace: {{ .CnvrgNs }}
spec:
  hosts:
    - "{{.Logging.Es.SvcName}}.{{ .ClusterDomain }}"
  gateways:
    - {{ .Networking.Istio.GwName }}
  http:
    - retries:
        attempts: {{ .Networking.Ingress.RetriesAttempts }}
        perTryTimeout: {{ .Networking.Ingress.PerTryTimeout }}
      timeout: {{ .Networking.Ingress.Timeout }}
      route:
        - destination:
            port:
              number: {{ .Logging.Es.Port}}
            host: "{{ .Logging.Es.SvcName }}.{{ .CnvrgNs }}.svc.cluster.local"
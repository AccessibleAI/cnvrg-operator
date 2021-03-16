apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Logging.Kibana.SvcName }}
  namespace: {{ .CnvrgNs }}
spec:
  hosts:
    - "{{.Logging.Kibana.SvcName}}.{{ .ClusterDomain }}"
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
          number: {{.Logging.Kibana.Port}}
        host: "{{ .Logging.Kibana.SvcName }}.{{ .CnvrgNs }}.svc.cluster.local"
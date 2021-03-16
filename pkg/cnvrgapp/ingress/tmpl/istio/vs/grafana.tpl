apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Monitoring.Grafana.SvcName }}
  namespace: {{ .CnvrgNs }}
spec:
  hosts:
    - "{{.Monitoring.Grafana.SvcName}}.{{ .ClusterDomain }}"
  gateways:
  - {{ .Networking.Istio.GwName }}
  http:
  - retries:
      attempts: {{ .Networking.Ingress.RetriesAttempts }}
      perTryTimeout: {{ .Networking.Ingress.PerTryTimeout }}
    timeout: {{ .Networking.Ingress.Timeout }}
    route:
    - destination:
        host: "{{ .Monitoring.Grafana.SvcName }}.{{ .CnvrgNs }}.svc.cluster.local"
        port:
          number: {{ .Monitoring.Grafana.Port }}

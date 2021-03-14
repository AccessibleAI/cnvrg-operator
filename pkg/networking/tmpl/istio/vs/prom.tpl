apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Monitoring.Prometheus.SvcName }}
  namespace: {{ .CnvrgNs }}
spec:
  hosts:
  - "{{ .Monitoring.Prometheus.SvcName }}.{{ .ClusterDomain }}"
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
          number: {{ .Monitoring.Prometheus.Port }}
        host: "{{ .Monitoring.Prometheus.SvcName }}.{{ .CnvrgNs }}.svc.cluster.local"
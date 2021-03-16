apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.Monitoring.Prometheus.SvcName }}
  namespace: {{ .Namespace }}
spec:
  hosts:
  - "{{ .Spec.Monitoring.Prometheus.SvcName }}.{{ .Spec.ClusterDomain }}"
  gateways:
  - {{ .Spec.Ingress.IstioGwName }}
  http:
  - retries:
      attempts: {{ .Spec.Ingress.RetriesAttempts }}
      perTryTimeout: {{ .Spec.Ingress.PerTryTimeout }}
    timeout: {{ .Spec.Ingress.Timeout }}
    route:
    - destination:
        port:
          number: {{ .Spec.Monitoring.Prometheus.Port }}
        host: "{{ .Spec.Monitoring.Prometheus.SvcName }}.{{ .Namespace }}.svc.cluster.local"
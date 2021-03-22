apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.Prometheus.SvcName }}
  namespace: {{ .Namespace }}
spec:
  hosts:
  - "{{ .Spec.Prometheus.SvcName }}.{{ .Spec.ClusterDomain }}"
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
          number: {{ .Spec.Prometheus.Port }}
        host: "{{ .Spec.Prometheus.SvcName }}.{{ .Namespace }}.svc.cluster.local"
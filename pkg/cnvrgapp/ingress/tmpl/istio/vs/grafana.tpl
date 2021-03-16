apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.Monitoring.Grafana.SvcName }}
  namespace: {{ .Namespace }}
spec:
  hosts:
    - "{{.Spec.Monitoring.Grafana.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
  - {{ .Spec.Ingress.IstioGwName }}
  http:
  - retries:
      attempts: {{ .Spec.Ingress.RetriesAttempts }}
      perTryTimeout: {{ .Spec.Ingress.PerTryTimeout }}
    timeout: {{ .Spec.Ingress.Timeout }}
    route:
    - destination:
        host: "{{ .Spec.Monitoring.Grafana.SvcName }}.{{ .Namespace }}.svc.cluster.local"
        port:
          number: {{ .Spec.Monitoring.Grafana.Port }}

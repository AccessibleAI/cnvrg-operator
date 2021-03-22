apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.Grafana.SvcName }}
  namespace: {{ ns . }}
spec:
  hosts:
    - "{{.Spec.Grafana.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
  - {{ .Spec.Ingress.IstioGwName }}
  http:
  - retries:
      attempts: {{ .Spec.Ingress.RetriesAttempts }}
      perTryTimeout: {{ .Spec.Ingress.PerTryTimeout }}
    timeout: {{ .Spec.Ingress.Timeout }}
    route:
    - destination:
        host: "{{ .Spec.Grafana.SvcName }}.{{ ns . }}.svc.cluster.local"
        port:
          number: {{ .Spec.Grafana.Port }}

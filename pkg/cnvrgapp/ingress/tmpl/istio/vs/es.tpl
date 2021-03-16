apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{.Spec.Logging.Es.SvcName }}
  namespace: {{ .Namespace }}
spec:
  hosts:
    - "{{.Spec.Logging.Es.SvcName}}.{{ .Spec.ClusterDomain }}"
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
              number: {{ .Spec.Logging.Es.Port}}
            host: "{{ .Spec.Logging.Es.SvcName }}.{{ .Namespace }}.svc.cluster.local"
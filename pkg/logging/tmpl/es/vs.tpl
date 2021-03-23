apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{.Spec.Logging.Es.SvcName }}
  namespace: {{ ns . }}
spec:
  hosts:
    - "{{.Spec.Logging.Es.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
    - {{ .Spec.Networking.Ingress.IstioGwName }}
  http:
    - retries:
        attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
        perTryTimeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
      timeout: {{ .Spec.Networking.Ingress.Timeout }}
      route:
        - destination:
            port:
              number: {{ .Spec.Logging.Es.Port}}
            host: "{{ .Spec.Logging.Es.SvcName }}.{{ ns . }}.svc.cluster.local"
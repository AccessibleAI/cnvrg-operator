apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
spec:
  hosts:
    - "{{.Spec.Logging.Kibana.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
  - {{ istioGwName .}}
  http:
  - retries:
      attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
      perTryTimeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
    timeout: {{ .Spec.Networking.Ingress.Timeout }}
    route:
    - destination:
        port:
          number: {{.Spec.Logging.Kibana.Port}}
        host: "{{ .Spec.Logging.Kibana.SvcName }}.{{ ns . }}.svc.cluster.local"
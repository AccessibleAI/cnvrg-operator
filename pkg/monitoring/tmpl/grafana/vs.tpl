apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.Monitoring.Grafana.SvcName }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
spec:
  hosts:
    - "{{.Spec.Monitoring.Grafana.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
  - {{ istioGwName .}}
  http:
  - retries:
      attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
      perTryTimeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
    timeout: {{ .Spec.Networking.Ingress.Timeout }}
    route:
    - destination:
        host: "{{ .Spec.Monitoring.Grafana.SvcName }}.{{ ns . }}.svc.cluster.local"
        port:
          number: {{ .Spec.Monitoring.Grafana.Port }}

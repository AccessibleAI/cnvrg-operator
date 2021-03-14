apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Minio.SvcName }}
  namespace: {{ .CnvrgNs }}
spec:
  hosts:
  - "{{ .Minio.SvcName}}.{{ .ClusterDomain }}"
  gateways:
  - {{ .Networking.Istio.GwName }}
  http:
  - retries:
      attempts: {{ .Networking.Ingress.RetriesAttempts }}
      perTryTimeout: 172800s
    timeout: 864000s
    route:
    - destination:
        host: "{{ .Minio.SvcName }}.{{ .CnvrgNs }}.svc.cluster.local"
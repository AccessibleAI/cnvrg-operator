apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.Minio.SvcName }}
  namespace: {{ .Namespace }}
spec:
  hosts:
  - "{{ .Spec.Minio.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
  - {{ .Spec.Ingress.IstioGwName }}
  http:
  - retries:
      attempts: {{ .Spec.Ingress.RetriesAttempts }}
      perTryTimeout: 172800s
    timeout: 864000s
    route:
    - destination:
        host: "{{ .Spec.Minio.SvcName }}.{{ .Namespace }}.svc.cluster.local"
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
spec:
  hosts:
  - "{{ .Spec.Dbs.Minio.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
  - {{ istioGwName .}}
  http:
  - retries:
      attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
      perTryTimeout: 172800s
    timeout: 864000s
    route:
    - destination:
        host: "{{ .Spec.Dbs.Minio.SvcName }}.{{ ns . }}.svc.cluster.local"
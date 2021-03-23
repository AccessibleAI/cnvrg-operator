apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.ControlPlan.Minio.SvcName }}
  namespace: {{ ns . }}
spec:
  hosts:
  - "{{ .Spec.ControlPlan.Minio.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
  - {{ .Spec.Networking.Ingress.IstioGwName }}
  http:
  - retries:
      attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
      perTryTimeout: 172800s
    timeout: 864000s
    route:
    - destination:
        host: "{{ .Spec.ControlPlan.Minio.SvcName }}.{{ ns . }}.svc.cluster.local"
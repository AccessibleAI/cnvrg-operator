apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.Dbs.Minio.SvcName }}
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
spec:
  hosts:
  - "{{ .Spec.Dbs.Minio.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
  - {{ .Spec.Networking.Ingress.IstioGwName}}
  http:
  - retries:
      attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
      perTryTimeout: 172800s
    timeout: 864000s
    route:
    - destination:
        host: "{{ .Spec.Dbs.Minio.SvcName }}.{{ ns . }}.svc.{{ .Spec.ClusterLocalDomain }}"
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{.Spec.Dbs.Es.SvcName }}
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
    - "{{.Spec.Dbs.Es.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
    - {{ .Spec.Networking.Ingress.IstioGwName}}
  http:
    - retries:
        attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
        perTryTimeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
      timeout: {{ .Spec.Networking.Ingress.Timeout }}
      route:
        - destination:
            port:
              number: {{ .Spec.Dbs.Es.Port}}
            host: "{{ .Spec.Dbs.Es.SvcName }}.{{ ns . }}.svc.{{ .Spec.ClusterInternalDomain }}"
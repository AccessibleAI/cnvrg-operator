apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: cnvrg-jwks
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
    - "{{ .Spec.Jwks.Name }}.{{ .Spec.ClusterDomain }}"
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
              number: 8080
            host: "{{ .Spec.Jwks.Name }}.{{ ns . }}.svc.{{ .Spec.ClusterInternalDomain }}"
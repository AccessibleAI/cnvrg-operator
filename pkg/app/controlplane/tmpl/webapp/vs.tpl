apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
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
    - "{{.Spec.ControlPlane.WebApp.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
    - {{ .Spec.Networking.Ingress.IstioGwName}}
  http:
    - retries:
        attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
        perTryTimeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
      timeout: {{ .Spec.Networking.Ingress.Timeout }}
      route:
        - destination:
            host: "{{ .Spec.ControlPlane.WebApp.SvcName }}.{{ ns . }}.svc.{{ .Spec.ClusterInternalDomain }}"
      headers:
        request:
          set:
            {{- if isTrue .Spec.Networking.HTTPS.Enabled }}
            x-forwarded-proto: https
            {{- else }}
            x-forwarded-proto: http
            {{- end}}


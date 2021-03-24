apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.ControlPlan.WebApp.SvcName }}
  namespace: {{ ns . }}
spec:
  hosts:
    - "{{.Spec.ControlPlan.WebApp.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
    - {{ .Spec.Networking.Ingress.IstioGwName }}
  http:
    - retries:
        attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
        perTryTimeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
      timeout: {{ .Spec.Networking.Ingress.Timeout }}
      route:
        - destination:
            host: "{{ .Spec.ControlPlan.WebApp.SvcName }}.{{ ns . }}.svc.cluster.local"
      headers:
        request:
          set:
            {{- if eq .Spec.Networking.HTTPS.Enabled "true"}}
            x-forwarded-proto: https
            {{- else }}
            x-forwarded-proto: http
            {{- end}}
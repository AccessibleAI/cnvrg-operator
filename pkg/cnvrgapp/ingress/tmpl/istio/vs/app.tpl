apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .ControlPlan.WebApp.SvcName }}
  namespace: {{ .CnvrgNs }}
spec:
  hosts:
    - "{{.ControlPlan.WebApp.SvcName}}.{{ .ClusterDomain }}"
  gateways:
    - {{ .Networking.Istio.GwName }}
  http:
    - retries:
        attempts: {{ .Networking.Ingress.RetriesAttempts }}
        perTryTimeout: {{ .Networking.Ingress.PerTryTimeout }}
      timeout: {{ .Networking.Ingress.Timeout }}
      route:
        - destination:
            host: "{{ .ControlPlan.WebApp.SvcName }}.{{ .CnvrgNs }}.svc.cluster.local"
      headers:
        request:
          set:
            {{- if eq .Networking.HTTPS.Enabled "true"}}
            x-forwarded-proto: https
            {{- else }}
            x-forwarded-proto: http
            {{- end}}
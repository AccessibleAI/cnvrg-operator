apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.ControlPlan.WebApp.SvcName }}
  namespace: {{ .Namespace }}
spec:
  hosts:
    - "{{.Spec.ControlPlan.WebApp.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
    - {{ .Spec.Ingress.IstioGwName }}
  http:
    - retries:
        attempts: {{ .Spec.Ingress.RetriesAttempts }}
        perTryTimeout: {{ .Spec.Ingress.PerTryTimeout }}
      timeout: {{ .Spec.Ingress.Timeout }}
      route:
        - destination:
            host: "{{ .Spec.ControlPlan.WebApp.SvcName }}.{{ .Namespace }}.svc.cluster.local"
      headers:
        request:
          set:
            {{- if eq .Spec.Ingress.HTTPS.Enabled "true"}}
            x-forwarded-proto: https
            {{- else }}
            x-forwarded-proto: http
            {{- end}}
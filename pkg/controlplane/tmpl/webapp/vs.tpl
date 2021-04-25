apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.ControlPlane.WebApp.SvcName }}
  namespace: {{ ns . }}
spec:
  hosts:
    - "{{.Spec.ControlPlane.WebApp.SvcName}}.{{ .Spec.ClusterDomain }}"
  gateways:
    - {{ istioGwName .}}
  http:
    - retries:
        attempts: {{ .Spec.Networking.Ingress.RetriesAttempts }}
        perTryTimeout: {{ .Spec.Networking.Ingress.PerTryTimeout }}
      timeout: {{ .Spec.Networking.Ingress.Timeout }}
      route:
        - destination:
            host: "{{ .Spec.ControlPlane.WebApp.SvcName }}.{{ ns . }}.svc.cluster.local"
      headers:
        request:
          set:
            {{- if .Spec.Networking.HTTPS.Enabled }}
            x-forwarded-proto: https
            {{- else }}
            x-forwarded-proto: http
            {{- end}}
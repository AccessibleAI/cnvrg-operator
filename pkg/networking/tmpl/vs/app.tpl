apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .ControlPlan.WebApp.SvcName }}
  namespace: {{ .CnvrgNs }}
spec:
  hosts:
    - {{ .ControlPlan.WebApp.SvcName}}.{{split ":" .ClusterDomain}}
  gateways:
    - {{ .Networking.Istio.GwName}}

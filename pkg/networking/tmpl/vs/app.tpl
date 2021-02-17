apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ .Spec.ControlPlan.WebApp.SvcName }}
  namespace: {{ .Spec.CnvrgNs }}
spec:
  hosts:
    - {{ .Spec.ControlPlan.WebApp.SvcName}}.{{split ":" .Spec.ClusterDomain}}
  gateways:
    - {{ .Spec.Networking.Istio.GwName}}

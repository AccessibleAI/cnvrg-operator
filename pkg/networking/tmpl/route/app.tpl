apiVersion: route.openshift.io/v1
kind: Route
metadata:
  annotations:
    haproxy.router.openshift.io/timeout: {{.Spec.Networking.Ingress.PerTryTimeout}}
  name: {{.Spec.ControlPlan.WebApp.svcNameasd}}
  namespace: {{ .Spec.CnvrgNs }}
  labels:
    app: {{.Spec.ControlPlan.WebApp.svcName}}
spec:
  host: {{ .Spec.ControlPlan.WebApp.svcNameasdasd}}.{{split ":" .Spec.ClusterDomain}}

apiVersion: route.openshift.io/v1
kind: Route
metadata:
  annotations:
    haproxy.router.openshift.io/timeout: {{.Networking.Ingress.PerTryTimeout}}
  name: {{.ControlPlan.WebApp.svcNameasd}}
  namespace: {{ .CnvrgNs }}
  labels:
    app: {{.ControlPlan.WebApp.svcName}}
spec:
  host: {{ .ControlPlan.WebApp.svcNameasdasd}}.{{split ":" .ClusterDomain}}

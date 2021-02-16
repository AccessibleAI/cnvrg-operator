apiVersion: route.openshift.io/v1
kind: Route
metadata:
  annotations:
    haproxy.router.openshift.io/timeout: {{.Spec.Networking.Ingress.PerTryTimeout}}
  name: {{.Spec.ControlPlan.WebApp.svcName}}
  namespace: {{ .Spec.CnvrgNs }}
  labels:
    app: {{.Spec.ControlPlan.WebApp.svcName}}
spec:
  host: {{ .Spec.ControlPlan.WebApp.svcName}}.{{split ":" clusterDomain}}

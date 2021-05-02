apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .Spec.Monitoring.Grafana.SvcName }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ .Spec.Monitoring.Grafana.SvcName }}
subjects:
  - kind: ServiceAccount
    name: {{ .Spec.Monitoring.Grafana.SvcName }}
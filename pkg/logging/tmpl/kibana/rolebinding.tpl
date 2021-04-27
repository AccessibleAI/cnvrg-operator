apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ .Spec.Logging.Kibana.SvcName }}
subjects:
  - kind: ServiceAccount
    name: {{ .Spec.Logging.Kibana.SvcName }}
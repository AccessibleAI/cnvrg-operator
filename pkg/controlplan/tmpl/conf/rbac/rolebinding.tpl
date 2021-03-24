apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .Spec.ControlPlan.Rbac.RoleBindingName }}
  namespace: {{ ns . }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ .Spec.ControlPlan.Rbac.Role }}
subjects:
  - kind: ServiceAccount
    name: {{ .Spec.ControlPlan.Rbac.ServiceAccountName }}
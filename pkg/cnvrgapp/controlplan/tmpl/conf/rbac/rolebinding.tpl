apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .ControlPlan.Rbac.RoleBindingName }}
  namespace: {{ .CnvrgNs }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ .ControlPlan.Rbac.Role }}
subjects:
  - kind: ServiceAccount
    name: {{ .ControlPlan.Rbac.ServiceAccountName }}
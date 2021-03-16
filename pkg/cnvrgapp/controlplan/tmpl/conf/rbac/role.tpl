apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ .ControlPlan.Rbac.Role }}
  namespace: {{ .CnvrgNs }}
rules:
  - apiGroups: ["*"]
    resources: ["*"]
    verbs: ["*"]
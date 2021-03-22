apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ .Spec.ControlPlan.Rbac.Role }}
  namespace: {{ ns . }}
rules:
  - apiGroups: ["*"]
    resources: ["*"]
    verbs: ["*"]
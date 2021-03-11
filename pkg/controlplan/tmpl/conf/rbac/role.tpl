apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cnvrg-control-plan
  namespace: {{ .CnvrgNs }}
rules:
  - apiGroups: ["*"]
    resources: ["*"]
    verbs: ["*"]
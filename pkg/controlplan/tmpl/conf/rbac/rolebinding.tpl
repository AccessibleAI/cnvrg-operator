apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: cnvrg-control-plan
  namespace: {{ .CnvrgNs }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cnvrg-control-plan
subjects:
  - kind: ServiceAccount
    name: cnvrg-control-plan
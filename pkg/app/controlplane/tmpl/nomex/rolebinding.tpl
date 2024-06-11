apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: "cnvrg-nomex"
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: "cnvrg-nomex"
subjects:
  - kind: ServiceAccount
    name: cnvrg-nomex
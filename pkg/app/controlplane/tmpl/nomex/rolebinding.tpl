apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: "cnvrg-nomex-{{ .Namespace }}"
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: "cnvrg-nomex-{{ .Namespace }}"
subjects:
  - kind: ServiceAccount
    name: cnvrg-nomex
    namespace: {{ .Namespace }}
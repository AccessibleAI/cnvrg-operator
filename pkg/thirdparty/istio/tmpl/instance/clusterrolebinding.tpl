kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: istio-operator
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
subjects:
  - kind: ServiceAccount
    name: istio-operator
    namespace: {{ .Namespace }}
roleRef:
  kind: ClusterRole
  name: istio-operator
  apiGroup: rbac.authorization.k8s.io
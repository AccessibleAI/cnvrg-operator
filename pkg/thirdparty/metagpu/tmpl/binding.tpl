apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: metagpu-device-plugin
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: metagpu-device-plugin
subjects:
  - kind: ServiceAccount
    name: metagpu-device-plugin
    namespace: {{ .Namespace }}
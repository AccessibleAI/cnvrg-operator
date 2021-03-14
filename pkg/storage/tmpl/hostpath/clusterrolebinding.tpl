apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: hostpath-provisioner
subjects:
  - kind: ServiceAccount
    name: hostpath-provisioner-admin
    namespace: {{ .CnvrgNs }}
roleRef:
  kind: ClusterRole
  name: hostpath-provisioner
  apiGroup: rbac.authorization.k8s.io
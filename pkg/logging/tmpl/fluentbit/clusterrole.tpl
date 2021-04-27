apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: fluent-bit-read
rules:
- apiGroups: [""]
  resources:
  - namespaces
  - pods
  verbs: ["get", "list", "watch"]
- apiGroups:
    - security.openshift.io
  resourceNames:
    - privileged
  resources:
    - securitycontextconstraints
  verbs:
    - use

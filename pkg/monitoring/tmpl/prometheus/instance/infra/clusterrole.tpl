apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cnvrg-infra-prometheus
  labels:
    owner: cnvrg-control-plane
rules:
- apiGroups:
    - security.openshift.io
  resourceNames:
    - anyuid
  resources:
    - securitycontextconstraints
  verbs:
    - use
- apiGroups:
    - ""
  resources:
    - services
    - pods
    - endpoints
  verbs:
    - get
    - list
    - watch
- apiGroups:
    - ""
  resources:
    - nodes/metrics
  verbs:
    - get
- nonResourceURLs:
    - /metrics
  verbs:
    - get
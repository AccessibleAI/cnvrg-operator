apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cnvrg-ccp-prometheus
  namespace: {{ ns . }}
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
    - endpoints
    - pods
  verbs:
    - get
    - list
    - watch
- apiGroups:
    - extensions
  resources:
    - ingresses
  verbs:
    - get
    - list
    - watch

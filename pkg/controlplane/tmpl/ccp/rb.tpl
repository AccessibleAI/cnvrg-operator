apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: cnvrg-ccp-operator-leader-election-rolebinding
  namespace: {{ ns . }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cnvrg-ccp-operator-leader-election-role
subjects:
  - kind: ServiceAccount
    name: cnvrg-ccp-operator-controller-manager
    namespace: {{ ns . }}
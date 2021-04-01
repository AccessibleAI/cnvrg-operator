kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: mpi-operator
  namespace: {{ ns . }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: mpi-operator
subjects:
  - kind: ServiceAccount
    name: mpi-operator
    namespace: {{ ns . }}



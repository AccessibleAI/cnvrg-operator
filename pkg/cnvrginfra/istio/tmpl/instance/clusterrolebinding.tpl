kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: istio-operator
subjects:
  - kind: ServiceAccount
    name: istio-operator
    namespace: {{ .Spec.CnvrgInfraNs }}
roleRef:
  kind: ClusterRole
  name: istio-operator
  apiGroup: rbac.authorization.k8s.io
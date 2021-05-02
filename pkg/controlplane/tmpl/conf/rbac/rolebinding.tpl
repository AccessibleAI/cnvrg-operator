apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .Spec.ControlPlane.Rbac.RoleBindingName }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ .Spec.ControlPlane.Rbac.Role }}
subjects:
  - kind: ServiceAccount
    name: {{ .Spec.ControlPlane.Rbac.ServiceAccountName }}
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .Spec.Logging.Elastalert.SvcName }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ .Spec.Logging.Elastalert.SvcName }}
subjects:
  - kind: ServiceAccount
    name: {{ .Spec.Logging.Elastalert.SvcName }}
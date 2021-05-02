apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .Spec.Dbs.Pg.ServiceAccount }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ .Spec.Dbs.Pg.ServiceAccount }}
subjects:
  - kind: ServiceAccount
    name: {{ .Spec.Dbs.Pg.ServiceAccount }}
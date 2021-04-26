apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ .Spec.Dbs.Es.ServiceAccount }}
  namespace: {{ ns . }}
rules:
- apiGroups:
  - "*"
  resources:
  - "*"
  verbs:
  - "*"
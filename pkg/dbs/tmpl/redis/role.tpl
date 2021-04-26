apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ .Spec.Dbs.Redis.ServiceAccount }}
  namespace: {{ ns . }}
rules:
- apiGroups:
  - "*"
  resources:
  - "*"
  verbs:
  - "*"
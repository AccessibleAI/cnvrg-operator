apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ .Spec.Dbs.Minio.ServiceAccount }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
rules:
- apiGroups:
  - "*"
  resources:
  - "*"
  verbs:
  - "*"
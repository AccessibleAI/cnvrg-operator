apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: cnvrg-ccp-prometheus
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cnvrg-ccp-prometheus
subjects:
- kind: ServiceAccount
  name: cnvrg-ccp-prometheus

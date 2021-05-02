apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cnvrg-infra-prometheus
  labels:
    owner: cnvrg-control-plane
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cnvrg-infra-prometheus
subjects:
- kind: ServiceAccount
  name: cnvrg-infra-prometheus
  namespace: {{ ns . }}

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cnvrg-prometheus-operator
  labels:
    app: cnvrg-prometheus-operator
    version: v0.44.1
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cnvrg-prometheus-operator
subjects:
- kind: ServiceAccount
  name: cnvrg-prometheus-operator
  namespace: {{ ns . }}

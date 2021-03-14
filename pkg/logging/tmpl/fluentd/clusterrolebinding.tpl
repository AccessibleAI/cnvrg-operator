kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: fluentd-clusterrole
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: fluentd-clusterrole
subjects:
- kind: ServiceAccount
  name: fluentd
  namespace: {{ .CnvrgNs }}

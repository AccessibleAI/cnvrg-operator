kind: ClusterRole 
apiVersion: rbac.authorization.k8s.io/v1
metadata: 
  name: fluentd-clusterrole 
rules: 
  - apiGroups: 
      - "" 
    resources: 
      - "namespaces" 
      - "pods" 
    verbs: 
      - "list"
      - "get" 
      - "watch"

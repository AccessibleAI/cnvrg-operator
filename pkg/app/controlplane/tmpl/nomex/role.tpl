apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cnvrg-nomex
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
rules:
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - list

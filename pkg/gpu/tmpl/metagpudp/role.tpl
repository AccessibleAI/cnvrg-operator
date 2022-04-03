apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: metagpu-device-plugin
  namespace: {{ .Namespace }}
rules:
  - apiGroups:
      - ""
    resources:
      - pods
    verbs:
      - list
      - get
  - apiGroups:
      - ""
    resources:
      - pods/exec
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - configmaps
    resourceNames:
      - metagpu-device-plugin-config
    verbs:
      - get
      - update
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  creationTimestamp: null
  name: cnvrg-ccp-operator-manager-role
  namespace: {{ ns . }}
rules:
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - mlops.cnvrg.io
    resources:
      - cnvrgclusterprovisioners
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - mlops.cnvrg.io
    resources:
      - cnvrgclusterprovisioners/finalizers
    verbs:
      - update
  - apiGroups:
      - mlops.cnvrg.io
    resources:
      - cnvrgclusterprovisioners/status
    verbs:
      - get
      - patch
      - update
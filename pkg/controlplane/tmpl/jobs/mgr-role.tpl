apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: workflowmap-manager-role
  namespace: {{ ns . }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
rules:
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - services
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
  - workflowmaps
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
  - workflowmaps/finalizers
  verbs:
  - update
- apiGroups:
  - mlops.cnvrg.io
  resources:
  - workflowmaps/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - networking.istio.io
  resources:
  - virtualservices
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cnvrg-jwks
  annotations:
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
  - secrets
  verbs:
  - watch
  - get
  - list
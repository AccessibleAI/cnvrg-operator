apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: prom
  namespace: {{ .Data.Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - services
  - endpoints
  verbs:
  - list
  - watch
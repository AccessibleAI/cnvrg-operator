apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: cnvrg-control-plane
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cnvrg-control-plane
subjects:
  - kind: ServiceAccount
    name: cnvrg-control-plane
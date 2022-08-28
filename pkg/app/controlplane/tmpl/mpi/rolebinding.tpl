kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: mpi-operator
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    owner: cnvrg-control-plane
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: mpi-operator
subjects:
  - kind: ServiceAccount
    name: mpi-operator
    namespace: {{ ns . }}



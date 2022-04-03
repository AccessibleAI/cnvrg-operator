apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: metagpu-device-plugin
  namespace: {{ .Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: metagpu-device-plugin
subjects:
  - kind: ServiceAccount
    name: metagpu-device-plugin
    namespace: {{ .Namespace }}
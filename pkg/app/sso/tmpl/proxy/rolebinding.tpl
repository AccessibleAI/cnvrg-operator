apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{.Spec.SSO.Proxy.SvcName}}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "true"
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
  name: {{.Spec.SSO.Proxy.SvcName}}
subjects:
  - kind: ServiceAccount
    name: {{.Spec.SSO.Proxy.SvcName}}

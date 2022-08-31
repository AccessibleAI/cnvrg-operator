apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .Spec.Dbs.Minio.ServiceAccount }}
  namespace: {{.Namespace }}
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
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ .Spec.Dbs.Minio.ServiceAccount }}
subjects:
  - kind: ServiceAccount
    name: {{ .Spec.Dbs.Minio.ServiceAccount }}
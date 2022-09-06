apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cnvrg-prom
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "false"
    mlops.cnvrg.io/updatable: "false"
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
  - list
  - watch
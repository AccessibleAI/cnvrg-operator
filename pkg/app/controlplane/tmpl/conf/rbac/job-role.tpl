apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cnvrg-job
  namespace: {{ .Namespace }}
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
{{- if not .Spec.ControlPlane.BaseConfig.CnvrgJobRbacStrict }}
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - services
  - configmaps
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - use
  - delete
{{- end }}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cnvrg-authz
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
  - "networking.k8s.io"
  - "networking.istio.io"
  - "route.openshift.io"
  - ""
  resources:
  - ingresses
  - virtualservices
  - routes
  - secrets
  verbs:
  - watch
  - get
  - list
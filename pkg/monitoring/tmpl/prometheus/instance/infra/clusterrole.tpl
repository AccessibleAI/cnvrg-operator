apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cnvrg-infra-prometheus
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
    - security.openshift.io
  resourceNames:
    - anyuid
  resources:
    - securitycontextconstraints
  verbs:
    - use
- apiGroups:
    - ""
  resources:
    - services
    - pods
    - endpoints
  verbs:
    - get
    - list
    - watch
- apiGroups:
    - ""
  resources:
    - nodes/metrics
  verbs:
    - get
- nonResourceURLs:
    - /metrics
  verbs:
    - get
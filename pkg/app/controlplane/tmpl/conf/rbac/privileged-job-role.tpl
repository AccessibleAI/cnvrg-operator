apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cnvrg-privileged-job
  namespace: {{ ns . }}
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
  - ""
  resources:
  - pods
  - services
  - configmaps
  - persistentvolumeclaims
  verbs:
  - '*'
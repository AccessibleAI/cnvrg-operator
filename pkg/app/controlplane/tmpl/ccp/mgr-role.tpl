apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cnvrg-ccp-operator-manager-role
  namespace: {{ ns . }}
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
rules:
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - mlops.cnvrg.io
    resources:
      - cnvrgclusterprovisioners
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - mlops.cnvrg.io
    resources:
      - cnvrgclusterprovisioners/finalizers
    verbs:
      - update
  - apiGroups:
      - mlops.cnvrg.io
    resources:
      - cnvrgclusterprovisioners/status
    verbs:
      - get
      - patch
      - update
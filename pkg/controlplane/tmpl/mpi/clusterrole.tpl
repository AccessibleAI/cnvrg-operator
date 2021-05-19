kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: mpi-operator
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    owner: cnvrg-control-plane
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
rules:
  - apiGroups:
      - ""
    resources:
      - configmaps
      - serviceaccounts
    verbs:
      - "*"
  # This is needed for the launcher Role.
  - apiGroups:
      - ""
    resources:
      - pods
    verbs:
      - get
      - list
      - watch
      - update
  # This is needed for the launcher Role.
  - apiGroups:
      - ""
    resources:
      - pods/exec
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - endpoints
    verbs:
      - create
      - get
      - update
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
  - apiGroups:
      - rbac.authorization.k8s.io
    resources:
      - roles
      - rolebindings
    verbs:
      - create
      - list
      - watch
  - apiGroups:
      - policy
    resources:
      - poddisruptionbudgets
    verbs:
      - create
      - list
      - update
      - watch
  - apiGroups:
      - apps
    resources:
      - statefulsets
    verbs:
      - create
      - list
      - update
      - watch
  - apiGroups:
      - batch
    resources:
      - jobs
    verbs:
      - create
      - list
      - update
      - watch
  - apiGroups:
      - apiextensions.k8s.io
    resources:
      - customresourcedefinitions
    verbs:
      - create
      - get
  - apiGroups:
      - kubeflow.org
    resources:
      - mpijobs
      - mpijobs/finalizers
      - mpijobs/status
    verbs:
      - "*"
  - apiGroups:
      - scheduling.incubator.k8s.io
      - scheduling.sigs.dev
      - scheduling.volcano.sh
    resources:
      - queues
      - podgroups
    verbs:
      - "*"
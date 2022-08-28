apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: istio-operator
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
      - authentication.istio.io
    resources:
      - '*'
    verbs:
      - '*'
  - apiGroups:
      - config.istio.io
    resources:
      - '*'
    verbs:
      - '*'
  - apiGroups:
      - install.istio.io
    resources:
      - '*'
    verbs:
      - '*'
  - apiGroups:
      - networking.istio.io
    resources:
      - '*'
    verbs:
      - '*'
  - apiGroups:
      - security.istio.io
    resources:
      - '*'
    verbs:
      - '*'
  # k8s groups
  - apiGroups:
      - admissionregistration.k8s.io
    resources:
      - mutatingwebhookconfigurations
      - validatingwebhookconfigurations
    verbs:
      - '*'
  - apiGroups:
      - apiextensions.k8s.io
    resources:
      - customresourcedefinitions.apiextensions.k8s.io
      - customresourcedefinitions
    verbs:
      - '*'
  - apiGroups:
      - apps
      - extensions
    resources:
      - daemonsets
      - deployments
      - deployments/finalizers
      - ingresses
      - replicasets
      - statefulsets
    verbs:
      - '*'
  - apiGroups:
      - autoscaling
    resources:
      - horizontalpodautoscalers
    verbs:
      - '*'
  - apiGroups:
      - monitoring.coreos.com
    resources:
      - servicemonitors
    verbs:
      - get
      - create
      - update
  - apiGroups:
      - policy
    resources:
      - poddisruptionbudgets
    verbs:
      - '*'
  - apiGroups:
      - rbac.authorization.k8s.io
    resources:
      - clusterrolebindings
      - clusterroles
      - roles
      - rolebindings
    verbs:
      - '*'
  - apiGroups:
      - ""
    resources:
      - configmaps
      - endpoints
      - events
      - namespaces
      - pods
      - persistentvolumeclaims
      - secrets
      - services
      - serviceaccounts
    verbs:
      - '*'
  - apiGroups:
      - coordination.k8s.io
    resources:
      - leases
    verbs:
      - '*'
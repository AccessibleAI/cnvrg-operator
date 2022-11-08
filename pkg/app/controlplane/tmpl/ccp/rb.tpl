apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: cnvrg-ccp-operator-leader-election-rolebinding
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
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cnvrg-ccp-operator-leader-election-role
subjects:
  - kind: ServiceAccount
    name: cnvrg-ccp-operator-controller-manager
    namespace: {{ ns . }}
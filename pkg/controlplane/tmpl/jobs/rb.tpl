apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: workflowmap-leader-election-rolebinding
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
  name: workflowmap-leader-election-role
subjects:
- kind: ServiceAccount
  name: workflowmap-controller-manager
  namespace: {{ ns . }}

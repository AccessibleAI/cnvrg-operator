apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ .Spec.Dbs.Es.Elastalert.SvcName }}
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
rules:
  - apiGroups:
      - security.openshift.io
    resourceNames:
      - anyuid
    resources:
      - securitycontextconstraints
    verbs:
      - use

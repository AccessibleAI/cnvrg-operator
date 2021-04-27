apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ .Spec.Logging.Kibana.SvcName }}
  namespace: {{ ns . }}
rules:
  - apiGroups:
      - security.openshift.io
    resourceNames:
      - anyuid
    resources:
      - securitycontextconstraints
    verbs:
      - use

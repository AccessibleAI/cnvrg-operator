apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Spec.Dbs.Es.ServiceAccount }}
  namespace: {{ ns . }}
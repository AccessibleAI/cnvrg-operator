apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Spec.Dbs.Redis.ServiceAccount }}
  namespace: {{ ns . }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
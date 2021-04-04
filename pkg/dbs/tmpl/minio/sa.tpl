apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Spec.Dbs.Minio.ServiceAccount }}
  namespace: {{ ns . }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Spec.Dbs.Pg.ServiceAccount }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
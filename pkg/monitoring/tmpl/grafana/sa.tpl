apiVersion: v1
kind: ServiceAccount
metadata:
  name: grafana
  namespace: {{ ns . }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
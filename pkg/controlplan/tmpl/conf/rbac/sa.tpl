apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Spec.ControlPlan.Rbac.ServiceAccountName }}
  namespace: {{ ns . }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
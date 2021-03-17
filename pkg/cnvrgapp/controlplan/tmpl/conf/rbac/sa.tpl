apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Spec.ControlPlan.Rbac.ServiceAccountName }}
  namespace: {{ .Namespace }}
imagePullSecrets:
  - name: {{ .Spec.ControlPlan.Registry.Name }}
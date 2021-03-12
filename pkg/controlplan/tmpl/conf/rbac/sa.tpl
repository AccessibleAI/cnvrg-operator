apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .ControlPlan.Rbac.ServiceAccountName }}
  namespace: {{ .CnvrgNs }}
imagePullSecrets:
  - name: {{ .ControlPlan.Registry.Name }}
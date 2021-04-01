apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Spec.ControlPlane.Rbac.ServiceAccountName }}
  namespace: {{ ns . }}
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
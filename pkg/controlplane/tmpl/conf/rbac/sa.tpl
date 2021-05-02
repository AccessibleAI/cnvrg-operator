apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Spec.ControlPlane.Rbac.ServiceAccountName }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
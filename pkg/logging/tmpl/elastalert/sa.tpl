apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Spec.Logging.Elastalert.SvcName }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
imagePullSecrets:
  - name: {{ .Spec.Registry.Name }}
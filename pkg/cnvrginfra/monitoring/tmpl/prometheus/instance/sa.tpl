apiVersion: v1
kind: ServiceAccount
metadata:
  name: cnvrg-infra-prometheus
  namespace: {{ .Spec.CnvrgInfraNs }}

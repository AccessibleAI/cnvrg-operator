apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: {{ .Spec.CnvrgInfraNs }}
  name: istio-operator
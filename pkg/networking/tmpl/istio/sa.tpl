apiVersion: v1
kind: ServiceAccount
metadata:
  namespace:  {{ .Spec.CnvrgNs }}
  name: istio-operator
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace:  {{ .CnvrgNs }}
  name: istio-operator
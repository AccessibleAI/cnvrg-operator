apiVersion: v1
kind: ServiceAccount
metadata:
  name: cnvrg-control-plan
  namespace: {{ .CnvrgNs }}
imagePullSecrets:
  - name: {{ .ControlPlan.Conf.Registry.Name }}
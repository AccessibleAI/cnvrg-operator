apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: cnvrg-istio
  namespace:  {{ .Spec.CnvrgNs }}
spec:
  profile: minimal
  namespace:  {{ .Spec.CnvrgNs }}
  hub: {{ .Spec.Networking.Istio.Hub }}
  tag: {{ .Spec.Networking.Istio.Tag }}
  values:
    global:
      istioNamespace:  {{ .Spec.CnvrgNs }}
      imagePullSecrets:
        - {{ .Spec.ControlPlan.Conf.Registry.Name }}
    meshConfig:
      rootNamespace:  {{ .Spec.CnvrgNs }}
  components:
    base:
      enabled: true
    pilot:
      enabled: true
      k8s:
        {{- if eq .Spec.ControlPlan.Conf.Tenancy.Enabled "true" }}
        nodeSelector:
          {{ .Spec.ControlPlan.Conf.Tenancy.Key }}: "{{ .Spec.ControlPlan.Conf.Tenancy.Value }}"
        {{- end }}
        tolerations:
        - key: "{{ .Spec.ControlPlan.Conf.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.ControlPlan.Conf.Tenancy.Value }}"
          effect: "NoSchedule"
    ingressGateways:
    - enabled: true
      name: istio-ingressgateway
      k8s:
        {{- if eq .Spec.ControlPlan.Conf.Tenancy.Enabled "true" }}
        nodeSelector:
          {{ .Spec.ControlPlan.Conf.Tenancy.Key }}: "{{ .Spec.ControlPlan.Conf.Tenancy.Value }}"
        {{- end }}
        tolerations:
        - key: "{{ .Spec.ControlPlan.Conf.Tenancy.Key }}"
          operator: "Equal"
          value: "{{ .Spec.ControlPlan.Conf.Tenancy.Value }}"
          effect: "NoSchedule"

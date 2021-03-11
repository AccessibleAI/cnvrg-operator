apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  name: "cnvrg-registry"
  namespace: {{ .CnvrgNs }}
data:
  .dockerconfigjson: {{ printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}" .ControlPlan.Conf.Registry.URL .ControlPlan.Conf.Registry.User .ControlPlan.Conf.Registry.Password  (printf "%s:%s" .ControlPlan.Conf.Registry.User .ControlPlan.Conf.Registry.Password | b64enc) | b64enc }}


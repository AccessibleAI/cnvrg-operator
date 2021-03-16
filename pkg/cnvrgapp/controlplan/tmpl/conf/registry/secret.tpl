apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  name: "cnvrg-registry"
  namespace: {{ .CnvrgNs }}
data:
  .dockerconfigjson: {{ printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}" .ControlPlan.Registry.URL .ControlPlan.Registry.User .ControlPlan.Registry.Password  (printf "%s:%s" .ControlPlan.Registry.User .ControlPlan.Registry.Password | b64enc) | b64enc }}


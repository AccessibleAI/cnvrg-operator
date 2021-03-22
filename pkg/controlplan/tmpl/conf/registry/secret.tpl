apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  name: "cnvrg-registry"
  namespace: {{ .Namespace }}
data:
  .dockerconfigjson: {{ printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}" .Spec.ControlPlan.Registry.URL .Spec.ControlPlan.Registry.User .Spec.ControlPlan.Registry.Password  (printf "%s:%s" .Spec.ControlPlan.Registry.User .Spec.ControlPlan.Registry.Password | b64enc) | b64enc }}


apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  name: {{ .Spec.ControlPlan.Mpi.Registry.Name }}
  namespace: {{ ns . }}
data:
  .dockerconfigjson: {{ printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}" .Spec.ControlPlan.Mpi.Registry.URL .Spec.ControlPlan.Mpi.Registry.User .Spec.ControlPlan.Mpi.Registry.Password  (printf "%s:%s" .Spec.ControlPlan.Mpi.Registry.User .Spec.ControlPlan.Mpi.Registry.Password | b64enc) | b64enc }}


apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  name: {{ .Spec.ControlPlane.Mpi.Registry.Name }}
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
data:
  .dockerconfigjson: {{ printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}" .Spec.ControlPlane.Mpi.Registry.URL .Spec.ControlPlane.Mpi.Registry.User .Spec.ControlPlane.Mpi.Registry.Password  (printf "%s:%s" .Spec.ControlPlane.Mpi.Registry.User .Spec.ControlPlane.Mpi.Registry.Password | b64enc) | b64enc }}


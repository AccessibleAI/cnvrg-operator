apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  name: {{ .Spec.ControlPlane.Mpi.Registry.Name }}
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    owner: cnvrg-control-plane
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  .dockerconfigjson: {{ printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}" .Spec.ControlPlane.Mpi.Registry.URL .Spec.ControlPlane.Mpi.Registry.User .Spec.ControlPlane.Mpi.Registry.Password  (printf "%s:%s" .Spec.ControlPlane.Mpi.Registry.User .Spec.ControlPlane.Mpi.Registry.Password | b64enc) | b64enc }}


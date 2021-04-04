apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  name: {{ .Data.Registry.Name }}
  namespace: {{ .Namespace }}
data:
  .dockerconfigjson: {{ printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}" .Data.Registry.URL .Data.Registry.User .Data.Registry.Password  (printf "%s:%s" .Data.Registry.User .Data.Registry.Password | b64enc) | b64enc }}


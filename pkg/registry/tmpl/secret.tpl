apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  name: "cnvrg-registry"
  namespace: {{ ns . }}
data:
  .dockerconfigjson: {{ printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}" .Spec.Registry.URL .Spec.Registry.User .Spec.Registry.Password  (printf "%s:%s" .Spec.Registry.User .Spec.Registry.Password | b64enc) | b64enc }}


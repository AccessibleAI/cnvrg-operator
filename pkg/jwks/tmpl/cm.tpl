apiVersion: v1
kind: ConfigMap
metadata:
  name: cnvrg-jwks
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  config.yaml: |-
    discovery:
      secret:
        labelKey: domainId
        dataKey: pub.key
    cache:
      enabled: true
      redis:
        address: localhost:6379
    api:
      listen: 0.0.0.0:8080
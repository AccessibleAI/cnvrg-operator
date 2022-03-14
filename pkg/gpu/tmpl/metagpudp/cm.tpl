apiVersion: v1
kind: ConfigMap
metadata:
  name: metagpu-device-plugin-config
  namespace: {{ .Namespace }}
  annotations:
    {{- range $k, $v := .Data.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    {{- range $k, $v := .Data.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  config.yaml: |
    accelerator: nvidia
    processesDiscoveryPeriod: 5
    deviceCacheTTL: 3600
    jwtSecret: topSecret
    mgctlTar: /tmp/mgctl
    mgctlAutoInject: true
    serverAddr: 0.0.0.0:50052
    memoryEnforcer: true
    deviceToken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1ldGFncHVAaW5zdGFuY2UiLCJ2aXNpYmlsaXR5TGV2ZWwiOiJsMCJ9.2rHykHFcHoIr-OCoPA5Am4ubf31-RJcayZnOTK6db94
    containerToken: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1ldGFncHVAaW5zdGFuY2UiLCJ2aXNpYmlsaXR5TGV2ZWwiOiJsMSJ9.o5v6Zdi1FKXQevRjuSbABBX1vIRYgN3Daz9iXabuFFA
    deviceSharing:
      - resourceName: cnvrg.io/metagpu
        autoReshare: true
        metaGpus: 2
        uuid: [ "*" ]
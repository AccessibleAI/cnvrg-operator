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
  {{- $secret := randAlphaNum 20 | b64enc }}
  {{- $deviceToken := generateMetagpuToken $secret "l0" }}
  {{- $containerToken := generateMetagpuToken $secret "l1" }}
  MG_EX_TOKEN: {{ $deviceToken }} # duplicated for the exporter
  MG_CTL_TOKEN: {{ $deviceToken }} # duplicated for the mgctl
  config.yaml: |
    accelerator: nvidia
    processesDiscoveryPeriod: 5
    deviceCacheTTL: 3600
    jwtSecret: {{ $secret }}
    mgctlTar: /tmp/mgctl
    mgctlAutoInject: true
    serverAddr: 0.0.0.0:50052
    memoryEnforcer: true
    deviceToken: {{ $deviceToken }}
    containerToken: {{ $containerToken }}
    deviceSharing:
      - resourceName: cnvrg.io/metagpu
        autoReshare: true
        metaGpus: 2
        uuid: [ "*" ]
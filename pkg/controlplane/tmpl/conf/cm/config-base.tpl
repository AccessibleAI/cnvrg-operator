apiVersion: v1
kind: ConfigMap
metadata:
  name: cp-base-config
  namespace: {{ ns . }}
  annotations:
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
    cnvrg-config-reloader.mlops.cnvrg.io: "autoreload-ccp"
    {{- range $k, $v := .Spec.Labels }}
    {{$k}}: "{{$v}}"
    {{- end }}
data:
  PASSENGER_APP_ENV: "app"
  RAILS_ENV: "app"
  RUN_JOBS_ON_SELF_CLUSTER: "true"
  DEFAULT_COMPUTE_NAME: "default"
  CHECK_JOB_EXPIRATION: "true"
  USE_STDOUT: "true"
  EXTRACT_TAGS_FROM_CMD: "true"
  AGENT_CUSTOM_TAG: "{{ .Spec.ControlPlane.BaseConfig.AgentCustomTag }}"
  KUBE_NAMESPACE: "{{ ns . }}"
  SHOW_INTERCOM: "{{ .Spec.ControlPlane.BaseConfig.Intercom }}"
  SPLIT_SIDEKIQ: "{{ .Spec.ControlPlane.Sidekiq.Split }}"
  CNVRG_PASSENGER_MAX_POOL_SIZE: "{{ .Spec.ControlPlane.WebApp.PassengerMaxPoolSize }}"
  OAUTH_PROXY_ENABLED: "{{ isTrue .Spec.SSO.Enabled }}"
  OAUTH_ADMIN_USER: "{{ .Spec.SSO.AdminUser }}"
  CNVRG_PASSENGER_BIND_ADDRESS: "{{ cnvrgPassengerBindAddress . }}"
  CNVRG_PASSENGER_BIND_PORT: "{{ cnvrgPassengerBindPort . }}"
  CNVRG_JOB_UID: "{{ .Spec.ControlPlane.BaseConfig.CnvrgJobUID }}"
  {{- if ne .Spec.ControlPlane.BaseConfig.JobsStorageClass "" }}
  CNVRG_JOBS_STORAGECLASS: "{{ .Spec.ControlPlane.BaseConfig.JobsStorageClass }}" # if is set, app's job will use this storageClass for notebooks/experiments
  {{- end }}
  {{- range $featureFlagName, $featureFlagValue := .Spec.ControlPlane.BaseConfig.FeatureFlags }}
  {{$featureFlagName}}: "{{$featureFlagValue}}"
  {{- end }}






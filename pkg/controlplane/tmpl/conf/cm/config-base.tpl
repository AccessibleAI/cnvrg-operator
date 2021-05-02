apiVersion: v1
kind: ConfigMap
metadata:
  name: cp-base-config
  namespace: {{ ns . }}
  labels:
    owner: cnvrg-control-plane
data:
  AGENT_CUSTOM_TAG: "{{ .Spec.ControlPlane.BaseConfig.AgentCustomTag }}"
  PASSENGER_APP_ENV: "{{ .Spec.ControlPlane.BaseConfig.PassengerAppEnv }}"
  RAILS_ENV: "{{ .Spec.ControlPlane.BaseConfig.RailsEnv }}"
  RUN_JOBS_ON_SELF_CLUSTER: "{{ .Spec.ControlPlane.BaseConfig.RunJobsOnSelfCluster }}"
  DEFAULT_COMPUTE_CONFIG: "{{ .Spec.ControlPlane.BaseConfig.DefaultComputeConfig }}"
  DEFAULT_COMPUTE_NAME: "{{ .Spec.ControlPlane.BaseConfig.DefaultComputeName }}"
  CHECK_JOB_EXPIRATION: "{{ .Spec.ControlPlane.BaseConfig.CheckJobExpiration }}"
  USE_STDOUT: "{{ .Spec.ControlPlane.BaseConfig.UseStdout }}"
  EXTRACT_TAGS_FROM_CMD: "{{ .Spec.ControlPlane.BaseConfig.ExtractTagsFromCmd }}"
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
  CNVRG_JOBS_STORAGECLASS: "{{ .Spec.ControlPlane.BaseConfig.JobsStorageClass }}" # if is set, app's job will use this storageClass for notebooeks/experiments
  {{- end }}
  {{- range $featureFlagName, $featureFlagValue := .Spec.ControlPlane.BaseConfig.FeatureFlags }}
  {{$featureFlagName}}: "{{$featureFlagValue}}"
  {{- end }}






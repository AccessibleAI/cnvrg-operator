apiVersion: v1
kind: ConfigMap
metadata:
  name: cp-base-config
  namespace: {{ .CnvrgNs }}
data:
  DEFAULT_COMPUTE_CLUSTER_DOMAIN: "{{ defaultComputeClusterDomain .}}"
  DEFAULT_COMPUTE_CLUSTER_HTTPS: "{{ .Networking.HTTPS.Enabled }}"
  AGENT_CUSTOM_TAG: "{{ .ControlPlan.BaseConfig.AgentCustomTag }}"
  PASSENGER_APP_ENV: "{{ .ControlPlan.BaseConfig.PassengerAppEnv }}"
  RAILS_ENV: "{{ .ControlPlan.BaseConfig.RailsEnv }}"
  RUN_JOBS_ON_SELF_CLUSTER: "{{ .ControlPlan.BaseConfig.RunJobsOnSelfCluster }}"
  DEFAULT_COMPUTE_CONFIG: "{{ .ControlPlan.BaseConfig.DefaultComputeConfig }}"
  DEFAULT_COMPUTE_NAME: "{{ .ControlPlan.BaseConfig.DefaultComputeName }}"
  CHECK_JOB_EXPIRATION: "{{ .ControlPlan.BaseConfig.CheckJobExpiration }}"
  USE_STDOUT: "{{ .ControlPlan.BaseConfig.UseStdout }}"
  EXTRACT_TAGS_FROM_CMD: "{{ .ControlPlan.BaseConfig.ExtractTagsFromCmd }}"
  KUBE_NAMESPACE: "{{ .CnvrgNs }}"
  SHOW_INTERCOM: "{{ .ControlPlan.BaseConfig.Intercom }}"
  SPLIT_SIDEKIQ: "{{ .ControlPlan.Sidekiq.Split }}"
  CNVRG_PASSENGER_MAX_POOL_SIZE: "{{ .ControlPlan.WebApp.PassengerMaxPoolSize }}"
  OAUTH_PROXY_ENABLED: "{{ .ControlPlan.OauthProxy.Enabled }}"
  OAUTH_ADMIN_USER: "{{ .ControlPlan.OauthProxy.AdminUser }}"
  CNVRG_PASSENGER_BIND_ADDRESS: "{{ cnvrgPassengerBindAddress . }}"
  CNVRG_PASSENGER_BIND_PORT: "{{ cnvrgPassengerBindPort . }}"
  CNVRG_JOB_UID: "{{ .ControlPlan.BaseConfig.CnvrgJobUID }}"
  {{- if ne .ControlPlan.BaseConfig.JobsStorageClass "" }}
  CNVRG_JOBS_STORAGECLASS: "{{ .ControlPlan.BaseConfig.JobsStorageClass }}" # if is set, app's job will use this storageClass for notebooeks/experiments
  {{- end }}
  {{- range $featureFlagName, $featureFlagValue := .ControlPlan.BaseConfig.FeatureFlags }}
  {{$featureFlagName}}: "{{$featureFlagValue}}"
  {{- end}}






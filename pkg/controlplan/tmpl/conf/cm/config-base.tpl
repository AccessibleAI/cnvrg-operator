apiVersion: v1
kind: ConfigMap
metadata:
  name: cp-base-config
  namespace: {{ .Namespace }}
data:
  DEFAULT_COMPUTE_CLUSTER_DOMAIN: "{{ defaultComputeClusterDomain .}}"
  DEFAULT_COMPUTE_CLUSTER_HTTPS: "{{ .Spec.Ingress.HTTPS.Enabled }}"
  AGENT_CUSTOM_TAG: "{{ .Spec.ControlPlan.BaseConfig.AgentCustomTag }}"
  PASSENGER_APP_ENV: "{{ .Spec.ControlPlan.BaseConfig.PassengerAppEnv }}"
  RAILS_ENV: "{{ .Spec.ControlPlan.BaseConfig.RailsEnv }}"
  RUN_JOBS_ON_SELF_CLUSTER: "{{ .Spec.ControlPlan.BaseConfig.RunJobsOnSelfCluster }}"
  DEFAULT_COMPUTE_CONFIG: "{{ .Spec.ControlPlan.BaseConfig.DefaultComputeConfig }}"
  DEFAULT_COMPUTE_NAME: "{{ .Spec.ControlPlan.BaseConfig.DefaultComputeName }}"
  CHECK_JOB_EXPIRATION: "{{ .Spec.ControlPlan.BaseConfig.CheckJobExpiration }}"
  USE_STDOUT: "{{ .Spec.ControlPlan.BaseConfig.UseStdout }}"
  EXTRACT_TAGS_FROM_CMD: "{{ .Spec.ControlPlan.BaseConfig.ExtractTagsFromCmd }}"
  KUBE_NAMESPACE: "{{ .Namespace }}"
  SHOW_INTERCOM: "{{ .Spec.ControlPlan.BaseConfig.Intercom }}"
  SPLIT_SIDEKIQ: "{{ .Spec.ControlPlan.Sidekiq.Split }}"
  CNVRG_PASSENGER_MAX_POOL_SIZE: "{{ .Spec.ControlPlan.WebApp.PassengerMaxPoolSize }}"
  OAUTH_PROXY_ENABLED: "{{ .Spec.ControlPlan.OauthProxy.Enabled }}"
  OAUTH_ADMIN_USER: "{{ .Spec.ControlPlan.OauthProxy.AdminUser }}"
  CNVRG_PASSENGER_BIND_ADDRESS: "{{ cnvrgPassengerBindAddress . }}"
  CNVRG_PASSENGER_BIND_PORT: "{{ cnvrgPassengerBindPort . }}"
  CNVRG_JOB_UID: "{{ .Spec.ControlPlan.BaseConfig.CnvrgJobUID }}"
  {{- if ne .Spec.ControlPlan.BaseConfig.JobsStorageClass "" }}
  CNVRG_JOBS_STORAGECLASS: "{{ .Spec.ControlPlan.BaseConfig.JobsStorageClass }}" # if is set, app's job will use this storageClass for notebooeks/experiments
  {{- end }}
  {{- range $featureFlagName, $featureFlagValue := .Spec.ControlPlan.BaseConfig.FeatureFlags }}
  {{$featureFlagName}}: "{{$featureFlagValue}}"
  {{- end}}






apiVersion: v1
kind: ConfigMap
metadata:
  name: cp-base-config
  namespace: {{ .CnvrgNs }}
data:
  DEFAULT_COMPUTE_CLUSTER_DOMAIN: {{ defaultComputeClusterDomain .}}
  DEFAULT_COMPUTE_CLUSTER_HTTPS: {{ .Networking.HTTPS.Enabled }}
  AGENT_CUSTOM_TAG: {{ .ControlPlan.Conf.AgentCustomTag }}
  PASSENGER_APP_ENV: {{ .ControlPlan.Conf.PassengerAppEnv }}
  RAILS_ENV: {{ .ControlPlan.Conf.RailsEnv }}
  RUN_JOBS_ON_SELF_CLUSTER: {{ .ControlPlan.Conf.RunJobsOnSelfCluster }}
  DEFAULT_COMPUTE_CONFIG {{ .ControlPlan.Conf.DefaultComputeConfig }}
  DEFAULT_COMPUTE_NAME: {{ .ControlPlan.Conf.DefaultComputeName }}
  CHECK_JOB_EXPIRATION: {{ .ControlPlan.Conf.CheckJobExpiration }}
  USE_STDOUT: {{ .ControlPlan.Conf.UseStdout }}
  EXTRACT_TAGS_FROM_CMD: {{ .ControlPlan.Conf.ExtractTagsFromCmd }}
  KUBE_NAMESPACE: {{ .CnvrgNs }}
  SHOW_INTERCOM: {{ .ControlPlan.Conf.Intercom }}
  SPLIT_SIDEKIQ: {{ .ControlPlan.Sidekiq.Split }}
  CNVRG_PASSENGER_MAX_POOL_SIZE: {{ .ControlPlan.WebApp.PassengerMaxPoolSize }}
  OAUTH_PROXY_ENABLED: {{ .ControlPlan.OauthProxy.Enabled }}
  OAUTH_ADMIN_USER: {{ .ControlPlan.OauthProxy.AdminUser }}
  CNVRG_PASSENGER_BIND_ADDRESS: {{ cnvrgPassengerBindAddress . }}
  CNVRG_PASSENGER_BIND_PORT: {{ cnvrgPassengerBindPort . }}
  CNVRG_JOB_UID: {{ .ControlPlan.Conf.CnvrgJobUID }}
  {{- if ne .ControlPlan.Conf.JobsStorageClass "" }}
  CNVRG_JOBS_STORAGECLASS: {{ .ControlPlan.Conf.JobsStorageClass }} # if is set, app's job will use this storageClass for notebooeks/experiments
  {{- end }}
  {{- range $featureFlagName, $featureFlagValue := .ControlPlan.Conf.FeatureFlags }}
  {{$featureFlagName}}: {{$featureFlagValue}}
  {{- end}}






apiVersion: v1
kind: ConfigMap
metadata:
  name: cp-base-config
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "true"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "true"
    {{- range $k, $v := .Spec.Annotations }}
    {{$k}}: "{{$v}}"
    {{- end }}
  labels:
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
  KUBE_NAMESPACE: "{{ .Namespace }}"
  SHOW_INTERCOM: "{{ .Spec.ControlPlane.BaseConfig.Intercom }}"
  SPLIT_SIDEKIQ: "{{ .Spec.ControlPlane.Sidekiq.Split }}"
  CNVRG_PASSENGER_MAX_POOL_SIZE: "{{ .Spec.ControlPlane.WebApp.PassengerMaxPoolSize }}"
  OAUTH_PROXY_ENABLED: "{{ isTrue .Spec.SSO.Enabled }}"
  OAUTH_ADMIN_USER: "{{ .Spec.SSO.Central.AdminUser }}"
  SSO_IDP_PRIVATE_KEY_REF: "{{ .Spec.SSO.Pki.PrivateKeySecret }}"
  CNVRG_PASSENGER_BIND_ADDRESS: "0.0.0.0"
  CNVRG_PASSENGER_BIND_PORT: "{{.Spec.ControlPlane.WebApp.Port}}"
  CNVRG_JOB_UID: "{{ .Spec.ControlPlane.BaseConfig.CnvrgJobUID }}"
  CNVRG_JOBS_SERVICE_ACCOUNT: "cnvrg-job"
  CNVRG_JOBS_PRIORITY_CLASS: "{{ .Spec.PriorityClass.JobClassRef }}"
  HAS_CNVRG_SCHEDULER: "{{ .Spec.ControlPlane.CnvrgScheduler.Enabled }}"
  CNVRG_SSO_SERVICE_URL: "{{ .Spec.SSO.Central.ServiceUrl }}"
  JWKS_ISS: "{{ .Spec.SSO.Central.JwksURL }}"
  JWKS_AUD: "cnvrg-tenant"
  CNVRG_BUILD_IMAGE_JOB_SA: "cnvrg-buildimage-job"
  CNVRG_ENDPOINTS_INDEX: "cnvrg-endpoints*"
  CNVRG_SPARK_JOB_SA: "cnvrg-spark-job"
  CVAT_ENABLED: "false"
  SSO_VERSION: "{{ .Spec.SSO.Version }}"
  CNVRG_RBAC_STRICT: "{{ .Spec.ControlPlane.BaseConfig.CnvrgJobRbacStrict }}"
  {{- if ne .Spec.ControlPlane.BaseConfig.JobsStorageClass "" }}
  CNVRG_JOBS_STORAGECLASS: "{{ .Spec.ControlPlane.BaseConfig.JobsStorageClass }}" # if is set, app's job will use this storageClass for notebooks/experiments
  {{- end }}
  {{- range $featureFlagName, $featureFlagValue := .Spec.ControlPlane.BaseConfig.FeatureFlags }}
  {{$featureFlagName}}: "{{$featureFlagValue}}"
  {{- end }}
  CNVRG_IMAGE_HUB: "{{ .Spec.ImageHub }}"






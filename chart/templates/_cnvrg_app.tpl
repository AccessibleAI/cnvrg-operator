{{- define "spec.cnvrgApp" }}
cnvrgApp:
  replicas: {{ .Values.cnvrgApp.replicas }}
  enabled: "{{ .Values.cnvrgApp.enabled }}"
  fixpg:  "{{ .Values.cnvrgApp.fixpg }}"
  image: "{{ .Values.cnvrgApp.image }}"
  port: "{{ .Values.cnvrgApp.port }}"
  svcName: "{{ .Values.cnvrgApp.svcName }}"
  nodePort: "{{ .Values.cnvrgApp.nodePort }}"
  passengerMaxPoolSize: {{ .Values.cnvrgApp.passengerMaxPoolSize }}
  enableReadinessProbe: "{{.Values.cnvrgApp.enableReadinessProbe}}"
  initialDelaySeconds: "{{ .Values.cnvrgApp.initialDelaySeconds }}"
  readinessPeriodSeconds: "{{ .Values.cnvrgApp.readinessPeriodSeconds }}"
  readinessTimeoutSeconds: "{{ .Values.cnvrgApp.readinessTimeoutSeconds }}"
  failureThreshold: "{{ .Values.cnvrgApp.failureThreshold }}"
  resourcesRequestEnabled: "{{.Values.cnvrgApp.resourcesRequestEnabled}}"

  {{- if eq .Values.computeProfile "large"}}
  cpu: "{{ .Values.computeProfiles.large.cnvrgApp.cpu }}"
  memory: "{{ .Values.computeProfiles.large.cnvrgApp.memory }}"
  sidekiq:
    enabled: "{{ .Values.cnvrgApp.sidekiq.enabled }}"
    split: "{{.Values.cnvrgApp.sidekiq.split }}"
    cpu: "{{ .Values.computeProfiles.large.cnvrgApp.sidekiq.cpu}}"
    memory: "{{ .Values.computeProfiles.large.cnvrgApp.sidekiq.memory}}"
    replicas: "{{ .Values.computeProfiles.large.cnvrgApp.sidekiq.replicas}}"
  searchkiq:
    enabled: "{{ .Values.cnvrgApp.searchkiq.enabled }}"
    cpu: "{{ .Values.computeProfiles.large.cnvrgApp.searchkiq.cpu}}"
    memory: "{{ .Values.computeProfiles.large.cnvrgApp.searchkiq.memory}}"
    replicas: "{{ .Values.computeProfiles.large.cnvrgApp.searchkiq.replicas}}"
  systemkiq:
    enabled: "{{ .Values.cnvrgApp.systemkiq.enabled }}"
    cpu: "{{ .Values.computeProfiles.large.cnvrgApp.systemkiq.cpu}}"
    memory: "{{ .Values.computeProfiles.large.cnvrgApp.systemkiq.memory}}"
    replicas: "{{ .Values.computeProfiles.large.cnvrgApp.systemkiq.replicas}}"
  {{- end }}

  {{- if eq .Values.computeProfile "medium"}}
  cpu: "{{ .Values.computeProfiles.medium.cnvrgApp.cpu }}"
  memory: "{{ .Values.computeProfiles.medium.cnvrgApp.memory }}"
  sidekiq:
    enabled: "{{ .Values.cnvrgApp.sidekiq.enabled }}"
    split: "{{.Values.cnvrgApp.sidekiq.split }}"
    cpu: "{{ .Values.computeProfiles.medium.cnvrgApp.sidekiq.cpu}}"
    memory: "{{ .Values.computeProfiles.medium.cnvrgApp.sidekiq.memory}}"
    replicas: "{{ .Values.computeProfiles.medium.cnvrgApp.sidekiq.replicas}}"
  searchkiq:
    enabled: "{{ .Values.cnvrgApp.searchkiq.enabled }}"
    cpu: "{{ .Values.computeProfiles.medium.cnvrgApp.searchkiq.cpu}}"
    memory: "{{ .Values.computeProfiles.medium.cnvrgApp.searchkiq.memory}}"
    replicas: "{{ .Values.computeProfiles.medium.cnvrgApp.searchkiq.replicas}}"
  systemkiq:
    enabled: "{{ .Values.cnvrgApp.systemkiq.enabled }}"
    cpu: "{{ .Values.computeProfiles.medium.cnvrgApp.systemkiq.cpu}}"
    memory: "{{ .Values.computeProfiles.medium.cnvrgApp.systemkiq.memory}}"
    replicas: "{{ .Values.computeProfiles.medium.cnvrgApp.systemkiq.replicas}}"
  {{- end }}

  {{- if eq .Values.computeProfile "small"}}
  cpu: "{{ .Values.computeProfiles.small.cnvrgApp.cpu }}"
  memory: "{{ .Values.computeProfiles.small.cnvrgApp.memory }}"
  sidekiq:
    enabled: "{{ .Values.cnvrgApp.sidekiq.enabled }}"
    split: "{{.Values.cnvrgApp.sidekiq.split }}"
    cpu: "{{ .Values.computeProfiles.small.cnvrgApp.sidekiq.cpu}}"
    memory: "{{ .Values.computeProfiles.small.cnvrgApp.sidekiq.memory}}"
    replicas: "{{ .Values.computeProfiles.small.cnvrgApp.sidekiq.replicas}}"
  searchkiq:
    enabled: "{{ .Values.cnvrgApp.searchkiq.enabled}}"
    cpu: "{{ .Values.computeProfiles.small.cnvrgApp.searchkiq.cpu}}"
    memory: "{{ .Values.computeProfiles.small.cnvrgApp.searchkiq.memory}}"
    replicas: "{{ .Values.computeProfiles.small.cnvrgApp.searchkiq.replicas}}"
  systemkiq:
    enabled: "{{ .Values.cnvrgApp.systemkiq.enabled}}"
    cpu: "{{ .Values.computeProfiles.small.cnvrgApp.systemkiq.cpu}}"
    memory: "{{ .Values.computeProfiles.small.cnvrgApp.systemkiq.memory}}"
    replicas: "{{ .Values.computeProfiles.small.cnvrgApp.systemkiq.replicas}}"
  {{- end }}
  seeder:
    image: "{{ .Values.cnvrgApp.seeder.image }}"
    seedCmd: "{{ .Values.cnvrgApp.seeder.seedCmd }}"
  conf:
    gcpStorageSecret: "{{ .Values.cnvrgApp.conf.gcpStorageSecret }}"
    gcpKeyfileMountPath: "{{ .Values.cnvrgApp.conf.gcpKeyfileMountPath }}"
    gcpKeyfileName: "{{ .Values.cnvrgApp.conf.gcpKeyfileName }}"
    jobsStorageClass: "{{ .Values.cnvrgApp.conf.jobsStorageClass }}"
    featureFlags: "{{ .Values.cnvrgApp.conf.featureFlags }}"
    sentryUrl: "{{ .Values.cnvrgApp.conf.sentryUrl }}"
    secretKeyBase: "{{ .Values.cnvrgApp.conf.secretKeyBase }}"
    stsIv: "{{ .Values.cnvrgApp.conf.stsIv }}"
    stsKey: "{{ .Values.cnvrgApp.conf.stsKey }}"
    passengerAppEnv: "{{ .Values.cnvrgApp.conf.passengerAppEnv }}"
    railsEnv: "{{ .Values.cnvrgApp.conf.railsEnv }}"
    runJobsOnSelfCluster: "{{ .Values.cnvrgApp.conf.runJobsOnSelfCluster }}"
    defaultComputeConfig: "{{ .Values.cnvrgApp.conf.defaultComputeConfig }}"
    defaultComputeName: "{{ .Values.cnvrgApp.conf.defaultComputeName }}"
    useStdout: "{{ .Values.cnvrgApp.conf.useStdout }}"
    extractTagsFromCmd: "{{ .Values.cnvrgApp.conf.extractTagsFromCmd }}"
    checkJobExpiration: "{{ .Values.cnvrgApp.conf.checkJobExpiration }}"
    cnvrgStorageType: "{{ .Values.cnvrgApp.conf.cnvrgStorageType }}"
    cnvrgStorageBucket: "{{ .Values.cnvrgApp.conf.cnvrgStorageBucket }}"
    cnvrgStorageAccessKey: "{{ .Values.cnvrgApp.conf.cnvrgStorageAccessKey }}"
    cnvrgStorageSecretKey: "{{ .Values.cnvrgApp.conf.cnvrgStorageSecretKey }}"
    {{- if ne .Values.cnvrgApp.conf.cnvrgStorageEndpoint "default" }}
    cnvrgStorageEndpoint: {{ .Values.cnvrgApp.conf.cnvrgStorageEndpoint }}
    {{- end}}
    minioSseMasterKey: "{{ .Values.cnvrgApp.conf.minioSseMasterKey }}"
    cnvrgStorageAzureAccessKey: "{{ .Values.cnvrgApp.conf.cnvrgStorageAzureAccessKey }}"
    cnvrgStorageAzureAccountName: "{{ .Values.cnvrgApp.conf.cnvrgStorageAzureAccountName }}"
    cnvrgStorageAzureContainer: "{{ .Values.cnvrgApp.conf.cnvrgStorageAzureContainer }}"
    cnvrgStorageRegion: "{{ .Values.cnvrgApp.conf.cnvrgStorageRegion }}"
    cnvrgStorageProject: "{{ .Values.cnvrgApp.conf.cnvrgStorageProject }}"
    customAgentTag: "{{ .Values.cnvrgApp.conf.customAgentTag }}"
    intercom: "{{ .Values.cnvrgApp.conf.intercom }}"
    ldap:
      enabled: "{{.Values.cnvrgApp.conf.ldap.enabled}}"
      host: "{{.Values.cnvrgApp.conf.ldap.host}}"
      port: "{{.Values.cnvrgApp.conf.ldap.port}}"
      account: "{{.Values.cnvrgApp.conf.ldap.account}}"
      base: "{{.Values.cnvrgApp.conf.ldap.base}}"
      adminUser: "{{.Values.cnvrgApp.conf.ldap.adminUser}}"
      adminPassword: "{{.Values.cnvrgApp.conf.ldap.adminPassword}}"
      ssl: "{{.Values.cnvrgApp.conf.ldap.ssl}}"
    registry:
      name: "{{ .Values.cnvrgApp.conf.registry.name}}"
      url: "{{ .Values.cnvrgApp.conf.registry.url}}"
      user: "{{ .Values.cnvrgApp.conf.registry.user}}"
      password: "{{ .Values.cnvrgApp.conf.registry.password}}"
    rbac:
      role: "{{ .Values.cnvrgApp.conf.rbac.role}}"
      serviceAccountName: "{{ .Values.cnvrgApp.conf.rbac.serviceAccountName}}"
      roleBindingName: "{{ .Values.cnvrgApp.conf.rbac.roleBindingName}}"
    smtp:
      server: "{{ .Values.cnvrgApp.conf.smtp.server}}"
      port: "{{ .Values.cnvrgApp.conf.smtp.port}}"
      username: "{{ .Values.cnvrgApp.conf.smtp.username}}"
      password: "{{ .Values.cnvrgApp.conf.smtp.password}}"
      domain: "{{ .Values.cnvrgApp.conf.smtp.domain}}"
  hyper:
    enabled: "{{ .Values.cnvrgApp.hyper.enabled }}"
    image: "{{ .Values.cnvrgApp.hyper.image }}"
    port: "{{ .Values.cnvrgApp.hyper.port }}"
    nodePort: "{{ .Values.cnvrgApp.hyper.nodePort }}"
    svcName: "{{ .Values.cnvrgApp.hyper.svcName }}"
    replicas: "{{ .Values.cnvrgApp.hyper.replicas }}"
    token: "{{.Values.cnvrgApp.hyper.token}}"
    enableReadinessProbe: "{{.Values.cnvrgApp.hyper.enableReadinessProbe}}"
    readinessPeriodSeconds: "{{.Values.cnvrgApp.hyper.readinessPeriodSeconds}}"
    readinessTimeoutSeconds: "{{.Values.cnvrgApp.hyper.readinessTimeoutSeconds}}"

    {{- if eq .Values.computeProfile "large"}}
    cpuRequest: "{{ .Values.computeProfiles.large.hyper.cpu }}"
    memoryRequest: "{{ .Values.computeProfiles.large.hyper.memory }}"
    {{- end }}

    {{- if eq .Values.computeProfile "medium"}}
    cpuRequest: "{{ .Values.computeProfiles.medium.hyper.cpu }}"
    memoryRequest: "{{ .Values.computeProfiles.medium.hyper.memory }}"
    {{- end }}

    {{- if eq .Values.computeProfile "small"}}
    cpuRequest: "{{ .Values.computeProfiles.small.hyper.cpu }}"
    memoryRequest: "{{ .Values.computeProfiles.small.hyper.memory }}"
    {{- end }}
  cnvrgRouter:
    enabled: "{{ .Values.cnvrgApp.cnvrgRouter.enabled }}"
    image: "{{.Values.cnvrgApp.cnvrgRouter.image}}"
    svcName: "{{ .Values.cnvrgApp.cnvrgRouter.svcName }}"
    nodePort: "{{ .Values.cnvrgApp.cnvrgRouter.nodePort }}"
    port: "{{ .Values.cnvrgApp.cnvrgRouter.port}}"

{{- end }}
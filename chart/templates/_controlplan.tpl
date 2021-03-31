{{- define "spec.controlPlan" }}
controlPlan:
  baseConfig:
    agentCustomTag: {{ .Values.controlPlan.baseConfig.agentCustomTag }}
    intercom: {{ .Values.controlPlan.baseConfig.intercom }}
    featureFlags: {{ .Values.controlPlan.baseConfig.featureFlags }}
  hyper:
    cpuLimit: {{ .Values.controlPlan.hyper.cpuLimit }}
    cpuRequest: {{ .Values.controlPlan.hyper.cpuRequest }}
    enableReadinessProbe: {{ .Values.controlPlan.hyper.enableReadinessProbe }}
    enabled: "{{ .Values.controlPlan.hyper.enabled }}"
    image: {{ .Values.controlPlan.hyper.image }}
    memoryLimit: {{ .Values.controlPlan.hyper.memoryLimit }}
    memoryRequest: {{ .Values.controlPlan.hyper.memoryRequest }}
    nodePort: {{ .Values.controlPlan.hyper.nodePort }}
    port: {{ .Values.controlPlan.hyper.port }}
    readinessPeriodSeconds: {{ .Values.controlPlan.hyper.readinessPeriodSeconds }}
    readinessTimeoutSeconds: {{ .Values.controlPlan.hyper.readinessTimeoutSeconds }}
    replicas: {{ .Values.controlPlan.hyper.replicas }}
    svcName: {{ .Values.controlPlan.hyper.svcName }}
    token: {{ .Values.controlPlan.hyper.token }}
  ldap:
    enabled: "{{ .Values.controlPlan.ldap.enabled }}"
    host: {{ .Values.controlPlan.ldap.host }}
    port: {{ .Values.controlPlan.ldap.port }}
    account: {{ .Values.controlPlan.ldap.account }}
    base: {{ .Values.controlPlan.ldap.base }}
    adminUser: {{ .Values.controlPlan.ldap.adminUser }}
    adminPassword: {{ .Values.controlPlan.ldap.adminPassword }}
    ssl: {{ .Values.controlPlan.ldap.ssl }}
  objectStorage:
    cnvrgStorageAccessKey: {{ .Values.controlPlan.objectStorage.cnvrgStorageAccessKey }}
    cnvrgStorageBucket: {{ .Values.controlPlan.objectStorage.cnvrgStorageBucket }}
    cnvrgStorageRegion: {{ .Values.controlPlan.objectStorage.cnvrgStorageRegion }}
    cnvrgStorageSecretKey: {{ .Values.controlPlan.objectStorage.cnvrgStorageSecretKey }}
    cnvrgStorageType: {{ .Values.controlPlan.objectStorage.cnvrgStorageType }}
    gcpKeyfileMountPath: {{ .Values.controlPlan.objectStorage.gcpKeyfileMountPath }}
    gcpKeyfileName: {{ .Values.controlPlan.objectStorage.gcpKeyfileName }}
    gcpStorageSecret: {{ .Values.controlPlan.objectStorage.gcpStorageSecret }}
    minioSseMasterKey: {{ .Values.controlPlan.objectStorage.minioSseMasterKey }}
    secretKeyBase: {{ .Values.controlPlan.objectStorage.secretKeyBase }}
    stsIv: {{ .Values.controlPlan.objectStorage.stsIv }}
    stsKey: {{ .Values.controlPlan.objectStorage.stsKey }}
  searchkiq:
    cpu: {{ .Values.controlPlan.searchkiq.cpu }}
    enabled: "{{ .Values.controlPlan.searchkiq.enabled }}"
    killTimeout: {{ .Values.controlPlan.searchkiq.killTimeout }}
    memory: {{ .Values.controlPlan.searchkiq.memory }}
    replicas: {{ .Values.controlPlan.searchkiq.replicas }}
  seeder:
    createBucketCmd: {{ .Values.controlPlan.seeder.createBucketCmd }}
    image: {{ .Values.controlPlan.seeder.image }}
    seedCmd: {{ .Values.controlPlan.seeder.seedCmd }}
  sidekiq:
    cpu: {{ .Values.controlPlan.sidekiq.cpu }}
    enabled: "{{ .Values.controlPlan.sidekiq.enabled }}"
    killTimeout: {{ .Values.controlPlan.sidekiq.killTimeout }}
    memory: {{ .Values.controlPlan.sidekiq.memory }}
    replicas: {{ .Values.controlPlan.sidekiq.replicas }}
    split: {{ .Values.controlPlan.sidekiq.split }}
  smtp:
    server: {{ .Values.controlPlan.smtp.server }}
    port: {{ .Values.controlPlan.smtp.port }}
    username: {{ .Values.controlPlan.smtp.username }}
    password: {{ .Values.controlPlan.smtp.password }}
    domain: {{ .Values.controlPlan.smtp.domain }}
  systemkiq:
    cpu: {{ .Values.controlPlan.systemkiq.cpu }}
    enabled: "{{ .Values.controlPlan.systemkiq.enabled }}"
    killTimeout: {{ .Values.controlPlan.systemkiq.killTimeout }}
    memory: {{ .Values.controlPlan.systemkiq.memory }}
    replicas: {{ .Values.controlPlan.systemkiq.replicas }}
  tenancy:
    dedicatedNodes: {{ .Values.controlPlan.tenancy.dedicatedNodes }}
    enabled: "{{ .Values.controlPlan.tenancy.enabled }}"
    key: {{ .Values.controlPlan.tenancy.key }}
    value: {{ .Values.controlPlan.tenancy.value }}
  webapp:
    cpu: {{ .Values.controlPlan.webapp.cpu }}
    enabled: "{{ .Values.controlPlan.webapp.enabled }}"
    failureThreshold: {{ .Values.controlPlan.webapp.failureThreshold }}
    image: {{ .Values.controlPlan.webapp.image }}
    initialDelaySeconds: {{ .Values.controlPlan.webapp.initialDelaySeconds }}
    memory: {{ .Values.controlPlan.webapp.memory }}
    nodePort: {{ .Values.controlPlan.webapp.nodePort }}
    oauthProxy:
      skipAuthRegex:
        - ^\/api
        - \/assets
        - \/healthz
        - \/public
        - \/pack
        - \/vscode.tar.gz
        - \/gitlens.vsix
        - \/ms-python-release.vsix
    passengerMaxPoolSize: {{ .Values.controlPlan.webapp.passengerMaxPoolSize }}
    port: {{ .Values.controlPlan.webapp.port }}
    readinessPeriodSeconds: {{ .Values.controlPlan.webapp.readinessPeriodSeconds }}
    readinessTimeoutSeconds: {{ .Values.controlPlan.webapp.readinessTimeoutSeconds }}
    replicas: {{ .Values.controlPlan.webapp.replicas }}
    svcName: {{ .Values.controlPlan.webapp.svcName }}
{{- end }}
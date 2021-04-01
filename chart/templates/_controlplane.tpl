{{- define "spec.controlPlane" }}
controlPlane:
  baseConfig:
    agentCustomTag: {{ .Values.controlPlane.baseConfig.agentCustomTag }}
    intercom: {{ .Values.controlPlane.baseConfig.intercom }}
    featureFlags: {{ .Values.controlPlane.baseConfig.featureFlags }}
  hyper:
    cpuLimit: {{ .Values.controlPlane.hyper.cpuLimit }}
    cpuRequest: {{ .Values.controlPlane.hyper.cpuRequest }}
    enableReadinessProbe: {{ .Values.controlPlane.hyper.enableReadinessProbe }}
    enabled: "{{ .Values.controlPlane.hyper.enabled }}"
    image: {{ .Values.controlPlane.hyper.image }}
    memoryLimit: {{ .Values.controlPlane.hyper.memoryLimit }}
    memoryRequest: {{ .Values.controlPlane.hyper.memoryRequest }}
    nodePort: {{ .Values.controlPlane.hyper.nodePort }}
    port: {{ .Values.controlPlane.hyper.port }}
    readinessPeriodSeconds: {{ .Values.controlPlane.hyper.readinessPeriodSeconds }}
    readinessTimeoutSeconds: {{ .Values.controlPlane.hyper.readinessTimeoutSeconds }}
    replicas: {{ .Values.controlPlane.hyper.replicas }}
    svcName: {{ .Values.controlPlane.hyper.svcName }}
    token: {{ .Values.controlPlane.hyper.token }}
  ldap:
    enabled: "{{ .Values.controlPlane.ldap.enabled }}"
    host: {{ .Values.controlPlane.ldap.host }}
    port: {{ .Values.controlPlane.ldap.port }}
    account: {{ .Values.controlPlane.ldap.account }}
    base: {{ .Values.controlPlane.ldap.base }}
    adminUser: {{ .Values.controlPlane.ldap.adminUser }}
    adminPassword: {{ .Values.controlPlane.ldap.adminPassword }}
    ssl: {{ .Values.controlPlane.ldap.ssl }}
  objectStorage:
    cnvrgStorageAccessKey: {{ .Values.controlPlane.objectStorage.cnvrgStorageAccessKey }}
    cnvrgStorageBucket: {{ .Values.controlPlane.objectStorage.cnvrgStorageBucket }}
    cnvrgStorageRegion: {{ .Values.controlPlane.objectStorage.cnvrgStorageRegion }}
    cnvrgStorageSecretKey: {{ .Values.controlPlane.objectStorage.cnvrgStorageSecretKey }}
    cnvrgStorageType: {{ .Values.controlPlane.objectStorage.cnvrgStorageType }}
    gcpKeyfileMountPath: {{ .Values.controlPlane.objectStorage.gcpKeyfileMountPath }}
    gcpKeyfileName: {{ .Values.controlPlane.objectStorage.gcpKeyfileName }}
    gcpStorageSecret: {{ .Values.controlPlane.objectStorage.gcpStorageSecret }}
    minioSseMasterKey: {{ .Values.controlPlane.objectStorage.minioSseMasterKey }}
    secretKeyBase: {{ .Values.controlPlane.objectStorage.secretKeyBase }}
    stsIv: {{ .Values.controlPlane.objectStorage.stsIv }}
    stsKey: {{ .Values.controlPlane.objectStorage.stsKey }}
  searchkiq:
    cpu: {{ .Values.controlPlane.searchkiq.cpu }}
    enabled: "{{ .Values.controlPlane.searchkiq.enabled }}"
    killTimeout: {{ .Values.controlPlane.searchkiq.killTimeout }}
    memory: {{ .Values.controlPlane.searchkiq.memory }}
    replicas: {{ .Values.controlPlane.searchkiq.replicas }}
  seeder:
    createBucketCmd: {{ .Values.controlPlane.seeder.createBucketCmd }}
    image: {{ .Values.controlPlane.seeder.image }}
    seedCmd: {{ .Values.controlPlane.seeder.seedCmd }}
  sidekiq:
    cpu: {{ .Values.controlPlane.sidekiq.cpu }}
    enabled: "{{ .Values.controlPlane.sidekiq.enabled }}"
    killTimeout: {{ .Values.controlPlane.sidekiq.killTimeout }}
    memory: {{ .Values.controlPlane.sidekiq.memory }}
    replicas: {{ .Values.controlPlane.sidekiq.replicas }}
    split: {{ .Values.controlPlane.sidekiq.split }}
  smtp:
    server: {{ .Values.controlPlane.smtp.server }}
    port: {{ .Values.controlPlane.smtp.port }}
    username: {{ .Values.controlPlane.smtp.username }}
    password: {{ .Values.controlPlane.smtp.password }}
    domain: {{ .Values.controlPlane.smtp.domain }}
  systemkiq:
    cpu: {{ .Values.controlPlane.systemkiq.cpu }}
    enabled: "{{ .Values.controlPlane.systemkiq.enabled }}"
    killTimeout: {{ .Values.controlPlane.systemkiq.killTimeout }}
    memory: {{ .Values.controlPlane.systemkiq.memory }}
    replicas: {{ .Values.controlPlane.systemkiq.replicas }}
  tenancy:
    dedicatedNodes: {{ .Values.controlPlane.tenancy.dedicatedNodes }}
    enabled: "{{ .Values.controlPlane.tenancy.enabled }}"
    key: {{ .Values.controlPlane.tenancy.key }}
    value: {{ .Values.controlPlane.tenancy.value }}
  webapp:
    cpu: {{ .Values.controlPlane.webapp.cpu }}
    enabled: "{{ .Values.controlPlane.webapp.enabled }}"
    failureThreshold: {{ .Values.controlPlane.webapp.failureThreshold }}
    image: {{ .Values.controlPlane.webapp.image }}
    initialDelaySeconds: {{ .Values.controlPlane.webapp.initialDelaySeconds }}
    memory: {{ .Values.controlPlane.webapp.memory }}
    nodePort: {{ .Values.controlPlane.webapp.nodePort }}
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
    passengerMaxPoolSize: {{ .Values.controlPlane.webapp.passengerMaxPoolSize }}
    port: {{ .Values.controlPlane.webapp.port }}
    readinessPeriodSeconds: {{ .Values.controlPlane.webapp.readinessPeriodSeconds }}
    readinessTimeoutSeconds: {{ .Values.controlPlane.webapp.readinessTimeoutSeconds }}
    replicas: {{ .Values.controlPlane.webapp.replicas }}
    svcName: {{ .Values.controlPlane.webapp.svcName }}
  mpi:
    enabled: "{{ .Values.controlPlane.mpi.enabled }}"
    image: {{ .Values.controlPlane.mpi.image }}
    kubectlDeliveryImage: {{ .Values.controlPlane.mpi.kubectlDeliveryImage }}
    extraArgs: {{ .Values.controlPlane.mpi.extraArgs }}
{{- end }}
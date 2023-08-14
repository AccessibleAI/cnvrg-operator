# cnvrg.io operator (v5)
---

## Deploy cnvrg stack on EKS | AKS | GKE | OpenShift | On-Premise clusters

## Add cnvrg helm chart repo

```bash
helm repo add cnvrg-operator https://charts.slim.cnvrg.io/cnvrg-operator
helm repo add cnvrg-cap https://charts.slim.cnvrg.io/cnvrg-cap
```

## Deploy cnvrg operator

```bash
helm install \
  cnvrg-operator \
  cnvrg-operator/cnvrg-operator \
  -n cnvrg-system --create-namespace \
  --set registry.user="<REGISTRY-USERNAME>" \
  --set registry.password="<REGISTRY-PASSWORD>"
```

## Deploy cnvrg instance

```bash
# simple deploy  
helm install \
  cnvrg \
  cnvrg-cap/cnvrg-cap \
  -n cnvrg --create-namespace \
  --set clusterDomain="<CLUSTER-DOMAIN>" \
  --set controlPlane.image="<CNVRG-APP-IMAGE>" \
  --set registry.user="<CNVRG-USERNAME>" \
  --set registry.password="<CNVRG-PASSWORD>"

# utilize non default ingress controller
# networking.ingress.type = ingress | istio | openshift  
helm install \
  cnvrg \
  cnvrg-cap/cnvrg-cap \
  -n cnvrg --create-namespace \
  --set clusterDomain="<CLUSTER-DOMAIN>" \
  --set controlPlane.image="<CNVRG-APP-IMAGE>" \
  --set registry.user="<CNVRG-USERNAME>" \
  --set registry.password="<CNVRG-PASSWORD>" \
  --set networking.ingress.type="ingress"
  

# deploy with SSO  
helm install \
  cnvrg \
  cnvrg-cap/cnvrg-cap \
  -n cnvrg --create-namespace \
  --set clusterDomain="<CLUSTER-DOMAIN>" \
  --set controlPlane.image="<CNVRG-APP-IMAGE>" \
  --set registry.user="<CNVRG-USERNAME>" \
  --set registry.password="<CNVRG-PASSWORD>" \
  --set sso.enabled=true \
  --set sso.jwks.enabled=true \
  --set sso.pki.enabled=true \
  --set sso.proxy.enabled=true \
  --set sso.central.enabled=true \
  --set sso.central.adminUser='<FIRST-ADMIN-USER>' \
  --set sso.central.provider='oidc' \
  --set sso.central.clientId='<CLIENT-ID>' \
  --set sso.central.clientSecret='<CLIENT-SECRET' \
  --set sso.central.oidcIssuerUrl='<OIDC-ISSUER-URL>'


# deploy on OpenShift
helm install \
  cnvrg \
  cnvrg-cap/cnvrg-cap \
  -n cnvrg --create-namespace \
  --set controlPlane.image="<CNVRG-APP-IMAGE>" \
  --set registry.user="<CNVRG-USERNAME>" \
  --set registry.password="<CNVRG-PASSWORD>" \
  --set networking.ingress.type="openshift" \
  --set networking.https.enabled="true"   
```

#### Helm chart deployment options

| **Flag**                                                           | **Default value**                   |
|--------------------------------------------------------------------|-------------------------------------|
| `clusterDomain`                                                    | -                                   |
| `clusterInternalDomain`                                            | cluster.local                       |
| `imageHub`                                                         | docker.io/cnvrg                     |
| `controlPlane.image`                                               | core:3.6.99                         |
| `controlPlane.webapp.replicas`                                     | 1                                   |
| `controlPlane.webapp.enabled`                                      | true                                |
| `controlPlane.webapp.port`                                         | 8080                                |
| `controlPlane.webapp.requests.cpu`                                 | 500m                                |
| `controlPlane.webapp.requests.memory`                              | 4Gi                                 |
| `controlPlane.webapp.limits.cpu`                                   | 4                                   |
| `controlPlane.webapp.limits.memory`                                | 8Gi                                 |
| `controlPlane.webapp.svcName`                                      | app                                 |
| `controlPlane.webapp.nodePort`                                     | 30080                               |
| `controlPlane.webapp.passengerMaxPoolSize`                         | 50                                  |
| `controlPlane.webapp.initialDelaySeconds`                          | 10                                  |
| `controlPlane.webapp.readinessPeriodSeconds`                       | 25                                  |
| `controlPlane.webapp.readinessTimeoutSeconds`                      | 20                                  |
| `controlPlane.webapp.failureThreshold`                             | 5                                   |
| `controlPlane.webapp.hpa.enabled`                                  | true                                |
| `controlPlane.webapp.hpa.utilization`                              | 85                                  |
| `controlPlane.webapp.hpa.maxReplicas`                              | 5                                   |
| `controlPlane.sidekiq.enabled`                                     | true                                |
| `controlPlane.sidekiq.split`                                       | true                                |
| `controlPlane.sidekiq.requests.cpu`                                | 200m                                |
| `controlPlane.sidekiq.requests.memory`                             | 3750Mi                              |
| `controlPlane.sidekiq.limits.cpu`                                  | 2                                   |
| `controlPlane.sidekiq.limits.memory`                               | 8Gi                                 |
| `controlPlane.sidekiq.replicas`                                    | 2                                   |
| `controlPlane.sidekiq.hpa.enabled`                                 | true                                |
| `controlPlane.sidekiq.hpa.utilization`                             | 85                                  |
| `controlPlane.sidekiq.hpa.maxReplicas`                             | 5                                   |
| `controlPlane.searchkiq.enabled`                                   | true                                |
| `controlPlane.searchkiq.requests.cpu`                              | 200m                                |
| `controlPlane.searchkiq.requests.memory`                           | 1Gi                                 |
| `controlPlane.searchkiq.limits.cpu`                                | 2                                   |
| `controlPlane.searchkiq.limits.memory`                             | 8Gi                                 |
| `controlPlane.searchkiq.replicas`                                  | 1                                   |
| `controlPlane.searchkiq.hpa.enabled`                               | true                                |
| `controlPlane.searchkiq.hpa.utilization`                           | 85                                  |
| `controlPlane.searchkiq.hpa.maxReplicas`                           | 5                                   |
| `controlPlane.systemkiq.enabled`                                   | true                                |
| `controlPlane.systemkiq.requests.cpu`                              | 300m                                |
| `controlPlane.systemkiq.requests.memory`                           | 2Gi                                 |
| `controlPlane.systemkiq.limits.cpu`                                | 2                                   |
| `controlPlane.systemkiq.limits.memory`                             | 8Gi                                 |
| `controlPlane.systemkiq.replicas`                                  | 1                                   |
| `controlPlane.systemkiq.hpa.enabled`                               | true                                |
| `controlPlane.systemkiq.hpa.utilization`                           | 85                                  |
| `controlPlane.systemkiq.hpa.maxReplicas`                           | 5                                   |
| `controlPlane.hyper.enabled`                                       | true                                |
| `controlPlane.hyper.image`                                         | hyper-server:latest                 |
| `controlPlane.hyper.port`                                          | 5050                                |
| `controlPlane.hyper.replicas`                                      | 1                                   |
| `controlPlane.hyper.nodePort`                                      | 30050                               |
| `controlPlane.hyper.svcName`                                       | hyper                               |
| `controlPlane.hyper.token`                                         | token                               |
| `controlPlane.hyper.requests.cpu`                                  | 100m                                |
| `controlPlane.hyper.requests.memory`                               | 200Mi                               |
| `controlPlane.hyper.limits.cpu`                                    | 2                                   |
| `controlPlane.hyper.limits.memory`                                 | 4Gi                                 |
| `controlPlane.hyper.cpuLimit`                                      | -                                   |
| `controlPlane.hyper.memoryLimit`                                   | -                                   |
| `controlPlane.hyper.readinessPeriodSeconds`                        | 100                                 |
| `controlPlane.hyper.readinessTimeoutSeconds`                       | 60                                  |
| `controlPlane.cnvrgScheduler.enabled`                              | false                               |
| `controlPlane.cnvrgScheduler.requests.cpu`                         | 200m                                |
| `controlPlane.cnvrgScheduler.requests.memory`                      | 1000Mi                              |
| `controlPlane.cnvrgScheduler.limits.cpu`                           | 2                                   |
| `controlPlane.cnvrgScheduler.limits.memory`                        | 4Gi                                 |
| `controlPlane.cnvrgScheduler.replicas`                             | 1                                   |
| `controlPlane.cnvrgRouter.enabled`                                 | false                               |
| `controlPlane.cnvrgRouter.image`                                   | nginx:1.21.0                        |
| `controlPlane.cnvrgRouter.svcName`                                 | cnvrg-router                        |
| `controlPlane.cnvrgRouter.nodePort`                                | 30081                               |
| `controlPlane.baseConfig.jobsStorageClass`                         | -                                   |
| `controlPlane.baseConfig.featureFlags.CNVRG_ENABLE_MOUNT_FOLDERS`  | false                               |
| `controlPlane.baseConfig.featureFlags.CNVRG_MOUNT_HOST_FOLDERS`    | false                               |
| `controlPlane.baseConfig.featureFlags.CNVRG_PROMETHEUS_METRICS`    | true                                |
| `controlPlane.baseConfig.sentryUrl`                                | -                                   |
| `controlPlane.baseConfig.runJobsOnSelfCluster`                     | -                                   |
| `controlPlane.baseConfig.agentCustomTag`                           | agnostic-logs                       |
| `controlPlane.baseConfig.intercom`                                 | true                                |
| `controlPlane.baseConfig.cnvrgJobUid`                              | 0                                   |
| `controlPlane.baseConfig.cnvrgJobRbacStrict`                       | false                               |
| `controlPlane.baseConfig.cnvrgPrivilegedJob`                       | true                                |
| `controlPlane.baseConfig.metagpuEnabled`                           | false                               |
| `controlPlane.ldap.enabled`                                        | false                               |
| `controlPlane.ldap.host`                                           | -                                   |
| `controlPlane.ldap.port`                                           | -                                   |
| `controlPlane.ldap.account`                                        | userPrincipalName                   |
| `controlPlane.ldap.base`                                           | -                                   |
| `controlPlane.ldap.adminUser`                                      | -                                   |
| `controlPlane.ldap.adminPassword`                                  | -                                   |
| `controlPlane.ldap.ssl`                                            | -                                   |
| `controlPlane.smtp.server`                                         | -                                   |
| `controlPlane.smtp.port`                                           | 587                                 |
| `controlPlane.smtp.username`                                       | -                                   |
| `controlPlane.smtp.password`                                       | -                                   |
| `controlPlane.smtp.domain`                                         | -                                   |
| `controlPlane.smtp.opensslVerifyMode`                              | -                                   |
| `controlPlane.smtp.sender`                                         | info@cnvrg.io                       |
| `controlPlane.objectStorage.type`                                  | minio                               |
| `controlPlane.objectStorage.bucket`                                | cnvrg-storage                       |
| `controlPlane.objectStorage.region`                                | eastus                              |
| `controlPlane.objectStorage.accessKey`                             | -                                   |
| `controlPlane.objectStorage.secretKey`                             | -                                   |
| `controlPlane.objectStorage.endpoint`                              | -                                   |
| `controlPlane.objectStorage.azureAccountName`                      | -                                   |
| `controlPlane.objectStorage.azureContainer`                        | -                                   |
| `controlPlane.objectStorage.gcpProject`                            | -                                   |
| `controlPlane.objectStorage.gcpSecretRef`                          | gcp-storage-secret                  |
| `controlPlane.mpi.enabled`                                         | false                               |
| `controlPlane.mpi.image`                                           | mpioperator/mpi-operator:v0.2.3     |
| `controlPlane.mpi.kubectlDeliveryImage`                            | mpioperator/kubectl-delivery:v0.2.3 |
| `controlPlane.mpi.extraArgs`                                       | null                                |
| `controlPlane.mpi.registry.name`                                   | mpi-private-registry                |
| `controlPlane.mpi.registry.url`                                    | docker.io                           |
| `controlPlane.mpi.registry.user`                                   | -                                   |
| `controlPlane.mpi.registry.password`                               | -                                   |
| `controlPlane.mpi.requests.cpu`                                    | 100m                                |
| `controlPlane.mpi.requests.memory`                                 | 100Mi                               |
| `controlPlane.mpi.limits.cpu`                                      | 1000m                               |
| `controlPlane.mpi.limits.memory`                                   | 1Gi                                 |
| `controlPlane.nomex.enabled`                                       | true                                |
| `controlPlane.nomex.image`                                         | nomex:v1.0.0                        |
| `registry.name`                                                    | cnvrg-app-registry                  |
| `registry.url`                                                     | docker.io                           |
| `registry.user`                                                    | -                                   |
| `registry.password`                                                | -                                   |
| `dbs.pg.enabled`                                                   | true                                |
| `dbs.pg.serviceAccount`                                            | pg                                  |
| `dbs.pg.image`                                                     | postgresql-12-centos7:latest        |
| `dbs.pg.port`                                                      | 5432                                |
| `dbs.pg.storageSize`                                               | 80Gi                                |
| `dbs.pg.svcName`                                                   | postgres                            |
| `dbs.pg.storageClass`                                              | -                                   |
| `dbs.pg.requests.cpu`                                              | 1                                   |
| `dbs.pg.requests.memory`                                           | 4Gi                                 |
| `dbs.pg.limits.cpu`                                                | 12                                  |
| `dbs.pg.limits.memory`                                             | 32Gi                                |
| `dbs.pg.maxConnections`                                            | 500                                 |
| `dbs.pg.sharedBuffers`                                             | 1024MB                              |
| `dbs.pg.effectiveCacheSize`                                        | 2048MB                              |
| `dbs.pg.hugePages.enabled`                                         | false                               |
| `dbs.pg.hugePages.size`                                            | 2Mi                                 |
| `dbs.pg.hugePages.memory`                                          | -                                   |
| `dbs.pg.nodeSelector`                                              | null                                |
| `dbs.pg.credsRef`                                                  | pg-creds                            |
| `dbs.pg.pvcName`                                                   | pg-storage                          |
| `dbs.redis.enabled`                                                | true                                |
| `dbs.redis.serviceAccount`                                         | redis                               |
| `dbs.redis.image`                                                  | cnvrg-redis:v3.0.5.c2               |
| `dbs.redis.svcName`                                                | redis                               |
| `dbs.redis.port`                                                   | 6379                                |
| `dbs.redis.storageSize`                                            | 10Gi                                |
| `dbs.redis.storageClass`                                           | -                                   |
| `dbs.redis.requests.cpu`                                           | 100m                                |
| `dbs.redis.requests.memory`                                        | 200Mi                               |
| `dbs.redis.limits.cpu`                                             | 1000m                               |
| `dbs.redis.limits.memory`                                          | 2Gi                                 |
| `dbs.redis.nodeSelector`                                           | null                                |
| `dbs.redis.credsRef`                                               | redis-creds                         |
| `dbs.redis.pvcName`                                                | redis-storage                       |
| `dbs.minio.enabled`                                                | true                                |
| `dbs.minio.serviceAccount`                                         | minio                               |
| `dbs.minio.replicas`                                               | 1                                   |
| `dbs.minio.image`                                                  | minio:RELEASE.2021-05-22T02-34-39Z  |
| `dbs.minio.port`                                                   | 9000                                |
| `dbs.minio.storageSize`                                            | 100Gi                               |
| `dbs.minio.svcName`                                                | minio                               |
| `dbs.minio.nodePort`                                               | 30090                               |
| `dbs.minio.storageClass`                                           | -                                   |
| `dbs.minio.requests.cpu`                                           | 200m                                |
| `dbs.minio.requests.memory`                                        | 2Gi                                 |
| `dbs.minio.limits.cpu`                                             | 8                                   |
| `dbs.minio.limits.memory`                                          | 20Gi                                |
| `dbs.minio.nodeSelector`                                           | null                                |
| `dbs.minio.pvcName`                                                | minio-storage                       |
| `dbs.es.enabled`                                                   | true                                |
| `dbs.es.serviceAccount`                                            | es                                  |
| `dbs.es.image`                                                     | cnvrg-es:7.17.5                     |
| `dbs.es.port`                                                      | 9200                                |
| `dbs.es.storageSize`                                               | 80Gi                                |
| `dbs.es.svcName`                                                   | elasticsearch                       |
| `dbs.es.nodePort`                                                  | 32200                               |
| `dbs.es.storageClass`                                              | -                                   |
| `dbs.es.requests.cpu`                                              | 500m                                |
| `dbs.es.requests.memory`                                           | 4Gi                                 |
| `dbs.es.limits.cpu`                                                | 4                                   |
| `dbs.es.limits.memory`                                             | 8Gi                                 |
| `dbs.es.javaOpts`                                                  | -                                   |
| `dbs.es.patchEsNodes`                                              | true                                |
| `dbs.es.nodeSelector`                                              | null                                |
| `dbs.es.credsRef`                                                  | es-creds                            |
| `dbs.es.pvcName`                                                   | es-storage                          |
| `dbs.es.cleanupPolicy.all`                                         | 3d                                  |
| `dbs.es.cleanupPolicy.app`                                         | 30d                                 |
| `dbs.es.cleanupPolicy.jobs`                                        | 14d                                 |
| `dbs.es.cleanupPolicy.endpoints`                                   | 1825d                               |
| `dbs.es.kibana.enabled`                                            | true                                |
| `dbs.es.kibana.serviceAccount`                                     | kibana                              |
| `dbs.es.kibana.svcName`                                            | kibana                              |
| `dbs.es.kibana.port`                                               | 8080                                |
| `dbs.es.kibana.image`                                              | kibana-oss:7.8.1                    |
| `dbs.es.kibana.nodePort`                                           | 30601                               |
| `dbs.es.kibana.requests.cpu`                                       | 100m                                |
| `dbs.es.kibana.requests.memory`                                    | 200Mi                               |
| `dbs.es.kibana.limits.cpu`                                         | 1000m                               |
| `dbs.es.kibana.limits.memory`                                      | 2Gi                                 |
| `dbs.es.kibana.credsRef`                                           | kibana-creds                        |
| `dbs.es.elastalert.enabled`                                        | true                                |
| `dbs.es.elastalert.image`                                          | elastalert:3.0.0-beta.1             |
| `dbs.es.elastalert.authProxyImage`                                 | nginx:1.20                          |
| `dbs.es.elastalert.credsRef`                                       | elastalert-creds                    |
| `dbs.es.elastalert.port`                                           | 8080                                |
| `dbs.es.elastalert.nodePort`                                       | 32030                               |
| `dbs.es.elastalert.storageSize`                                    | 30Gi                                |
| `dbs.es.elastalert.svcName`                                        | elastalert                          |
| `dbs.es.elastalert.storageClass`                                   | -                                   |
| `dbs.es.elastalert.requests.cpu`                                   | 100m                                |
| `dbs.es.elastalert.requests.memory`                                | 200Mi                               |
| `dbs.es.elastalert.limits.cpu`                                     | 400m                                |
| `dbs.es.elastalert.limits.memory`                                  | 800Mi                               |
| `dbs.es.elastalert.nodeSelector`                                   | null                                |
| `dbs.es.elastalert.pvcName`                                        | elastalert-storage                  |
| `dbs.prom.enabled`                                                 | true                                |
| `dbs.prom.credsRef`                                                | prom-creds                          |
| `dbs.prom.extraScrapeConfigs`                                      | null                                |
| `dbs.prom.image`                                                   | prometheus:v2.37.1                  |
| `dbs.prom.storageClass`                                            | ""                                  |
| `dbs.prom.storageSize`                                             | 50Gi                                |
| `dbs.prom.grafana.enabled`                                         | true                                |
| `dbs.prom.grafana.image`                                           | grafana-oss:9.1.7                   |
| `dbs.prom.grafana.svcName`                                         | grafana                             |
| `dbs.prom.grafana.port`                                            | 8080                                |
| `dbs.prom.grafana.nodePort`                                        | 30012                               |
| `dbs.prom.grafana.credsRef`                                        | grafana-creds                       |
| `networking.ingress.type`                                          | istio                               |
| `networking.ingress.timeout`                                       | 18000s                              |
| `networking.ingress.retriesAttempts`                               | 5                                   |
| `networking.ingress.perTryTimeout`                                 | 3600s                               |
| `networking.ingress.istioGwEnabled`                                | true                                |
| `networking.ingress.istioGwName`                                   | -                                   |
| `networking.ingress.istioIngressSelectorKey`                       | istio                               |
| `networking.ingress.istioIngressSelectorValue`                     | ingressgateway                      |
| `networking.ingress.ocpSecureRoutes`                               | false                               |
| `networking.https.enabled`                                         | false                               |
| `networking.https.certSecret`                                      | -                                   |
| `networking.https.cert`                                            | -                                   |
| `networking.https.key`                                             | -                                   |
| `networking.proxy.enabled`                                         | false                               |
| `networking.proxy.configRef`                                       | cp-proxy                            |
| `sso.enabled`                                                      | false                               |
| `sso.version`                                                      | v3                                  |
| `sso.pki.enabled`                                                  | false                               |
| `sso.pki.rootCaSecret`                                             | sso-idp-root-ca                     |
| `sso.pki.privateKeySecret`                                         | sso-idp-private-key                 |
| `sso.pki.publicKeySecret`                                          | sso-idp-pki-public-key              |
| `sso.jwks.enabled`                                                 | false                               |
| `sso.jwks.name`                                                    | cnvrg-jwks                          |
| `sso.jwks.image`                                                   | jwks:latest                         |
| `sso.jwks.cacheImage`                                              | redis:7.0.5                         |
| `sso.central.enabled`                                              | false                               |
| `sso.central.publicUrl`                                            | -                                   |
| `sso.central.oauthProxyImage`                                      | oauth2-proxy:v7.4.ssov3.p6          |
| `sso.central.centralUiImage`                                       | centralsso:latest                   |
| `sso.central.adminUser`                                            | -                                   |
| `sso.central.provider`                                             | -                                   |
| `sso.central.emailDomain.0`                                        | *                                   |
| `sso.central.clientId`                                             | -                                   |
| `sso.central.clientSecret`                                         | -                                   |
| `sso.central.oidcIssuerUrl`                                        | -                                   |
| `sso.central.serviceUrl`                                           | -                                   |
| `sso.central.scope`                                                | openid email profile                |
| `sso.central.insecureOidcAllowUnverifiedEmail`                     | true                                |
| `sso.central.whitelistDomain`                                      | -                                   |
| `sso.central.cookieDomain`                                         | -                                   |
| `sso.central.groupsAuth`                                           | false                               |
| `sso.central.readiness`                                            | true                                |
| `sso.central.requests.cpu`                                         | 500m                                |
| `sso.central.requests.memory`                                      | 1Gi                                 |
| `sso.central.limits.cpu`                                           | 2                                   |
| `sso.central.limits.memory`                                        | 4Gi                                 |
| `sso.proxy.enabled`                                                | false                               |
| `sso.proxy.image`                                                  | cnvrg-proxy:v1.0.15                 |
| `sso.proxy.address`                                                | -                                   |
| `sso.proxy.readiness`                                              | true                                |
| `sso.proxy.requests.cpu`                                           | 500m                                |
| `sso.proxy.requests.memory`                                        | 1Gi                                 |
| `sso.proxy.limits.cpu`                                             | 2                                   |
| `sso.proxy.limits.memory`                                          | 4Gi                                 |
| `tenancy.enabled`                                                  | false                               |
| `tenancy.key`                                                      | purpose                             |
| `tenancy.value`                                                    | cnvrg-control-plane                 |
| `priorityClass.appClassRef`                                        | -                                   |
| `priorityClass.jobClassRef`                                        | -                                   |
# cnvrg.io operator (v3)
---
## Deploy cnvrg stack on EKS | AKS | GKE | OpenShift | On-Premise clusters

### Architecture overview 
cnvrg operator may deploy cnvrg stack in two different ways
1. Multiple cnvrg control plans within the same cluster separated by namespaces - suitable for multi tenancy deployments  
```shell
                            ---------cnvrg infra namespace----------
                            | Cluster scope prometheus             |
                            | Prometheus node exporter             |
                            | Kube state metrics                   |
                            | Cluster scope service monitors       |     
                            | Fluentbit                            |
                            | Istio control plan                   |
                            | Storage provisioners (hostpath/nfs)  |
                            ----------------------------------------           
---------cnvrg control plan 1 namespace-------  ---------cnvrg control plan 2 namespace-------
| cnvrg control plan (webapp, sidekiqs, etc.)|  | cnvrg control plan (webapp, sidekiqs, etc.)|
| PostgreSQL                                 |  | PostgreSQL                                 |
| ElasticSearch + Kibana                     |  | ElasticSearch + Kibana                     |
| Minio                                      |  | Minio                                      |
| Redis                                      |  | Redis                                      |
| Namespace scope Prometheus + Grafana       |  | Namespace scope Prometheus + Grafana       |
| Namespace scope service monitors           |  | Namespace scope service monitors           |
| Istio Gateway + VirtualServices            |  | Istio Gateway + VirtualServices            |
----------------------------------------------  ----------------------------------------------
                    
```
2. Single cnvrg control plan in dedicated namespace 
```shell
                        ----------------cnvrg namespace--------------------
                        | Cluster scope prometheus                        |
                        | Prometheus node exporter                        |
                        | Kube state metrics                              |
                        | Cluster scope service monitors                  |     
                        | Namespace scope service monitors                |     
                        | Fluentbit                                       |
                        | Istio control plan                              |
                        | Storage provisioners (hostpath/nfs)             |   
                        | cnvrg control plan (webapp, sidekiqs, etc.)     |
                        | PostgreSQL                                      |
                        | ElasticSearch + Kibana                          | 
                        | Minio                                           |
                        | Redis                                           |  
                        | IstioGateway + VirtualServices                  |
                        ---------------------------------------------------           
```



### Quick start - namespace tenancy `disabled` with single cnvrg control plane
```shell
helm install cnvrg-1 . -n cnvrg-1 --create-namespace \ 
  --set appClusterDomain="<control-plane-domain-record>" \
  --set controlPlane.webapp.image="<cnvrg-control-plane-image>" \
  --set registry.user="<cnvrg-private-registry-user>" \
  --set registry.password="<cnvrg-private-registry-password>"  
```

### Quick start - namespace tenancy `enabled` with multiple cnvrg control planes within the same K8s cluster 

Deploy cnvrg infrastructure first 
```shell
helm install cnvrg-infra . -n cnvrg-infra --create-namespace \
  --set namespaceTenancy="true" \
  --set infraClusterDomain="<infrastructure-domain-record>" \
  --set registry.user="<cnvrg-private-registry-user>" \
  --set registry.password="<cnvrg-private-registry-password>"  
```
Once infrastructure components are ready, deploy cnvrg control plane 
```shell
helm install cnvrg-1 . -n cnvrg-1 --create-namespace \
  --set namespaceTenancy="true" \
  --set appClusterDomain="<control-plane-domain-record>" \
  --set controlPlane.webapp.image="<cnvrg-control-plane-image>" \
  --set registry.user="<cnvrg-private-registry-user>" \
  --set registry.password="<cnvrg-private-registry-password>"  
```


### Examples 
enable on-prem nfs storage  
```shell
  ... 
  --set storage.enabled="true" \
  --set storage.nfs.enabled="true" \
  --set storage.nfs.defaultSc="true" \
  --set storage.nfs.server="<NFS-SERVER-IP>" \
  --set storage.nfs.path="<EXPORT-PATH>" \
  ...  
```

enable on-prom istio with TLS termination  
```shell
  ... 
  --set networking.https.enabled="true" \
  --set networking.https.certSecret="<CERTIFICATE-K8S-SECRET>" \
  --set networking.istio.externalIp="<K8S-NODES-IPS>" \
  ...  
```

enable SSO 
```shell
  ... 
  --set sso.enabled="true" \
  --set sso.adminUser="<admin user>" \
  --set sso.provider="<provider>" \
  --set sso.emailDomain="<email-domain>" \
  --set sso.clientId="<client-id>" \
  --set sso.clientSecret="<client-secret>" \ 
  ...  
```

#### External Monitoring with OpenShift
1. Get user and password for OpenShift Prometheus instance 
```bash
CNVRG_PROMETHEUS_USER=$(kubectl get secret -nopenshift-monitoring grafana-datasources -ojson | jq -r '.data."prometheus.yaml"' | base64 -D | jq -r '.datasources[].basicAuthUser')
CNVRG_PROMETHEUS_PASS=$(kubectl get secret -nopenshift-monitoring grafana-datasources -ojson | jq -r '.data."prometheus.yaml"' | base64 -D | jq -r '.datasources[].basicAuthPassword')
CNVRG_PROMETHEUS_URL=$(kubectl get secret -nopenshift-monitoring grafana-datasources -ojson | jq -r '.data."prometheus.yaml"' | base64 -D | jq -r '.datasources[].url')

kubectl create secret generic prom-creds -ncnvrg \
 --from-literal=CNVRG_PROMETHEUS_USER=$CNVRG_PROMETHEUS_USER \
 --from-literal=CNVRG_PROMETHEUS_PASS=$CNVRG_PROMETHEUS_PASS \
 --from-literal=CNVRG_PROMETHEUS_URL=$CNVRG_PROMETHEUS_URL

```

### Chart options 
|**key**|**default value**
| ---|---| 
|`clusterDomain`|-
|`controlPlane.baseConfig.agentCustomTag`|latest
|`controlPlane.baseConfig.ccpStorageClass`|-
|`controlPlane.baseConfig.checkJobExpiration`|true
|`controlPlane.baseConfig.cnvrgJobUid`|1000
|`controlPlane.baseConfig.defaultComputeConfig`|/opt/kube
|`controlPlane.baseConfig.defaultComputeName`|default
|`controlPlane.baseConfig.extractTagsFromCmd`|false
|`controlPlane.baseConfig.hostpathNode`|-
|`controlPlane.baseConfig.intercom`|true
|`controlPlane.baseConfig.jobsStorageClass`|-
|`controlPlane.baseConfig.passengerAppEnv`|app
|`controlPlane.baseConfig.railsEnv`|app
|`controlPlane.baseConfig.runJobsOnSelfCluster`|true
|`controlPlane.baseConfig.useStdout`|true
|`controlPlane.hyper.cpuLimit`|2000m
|`controlPlane.hyper.cpuRequest`|100m
|`controlPlane.hyper.enableReadinessProbe`|true
|`controlPlane.hyper.enabled`|true
|`controlPlane.hyper.image`|cnvrg/hyper-server:latest
|`controlPlane.hyper.memoryLimit`|4Gi
|`controlPlane.hyper.memoryRequest`|200Mi
|`controlPlane.hyper.nodePort`|30050
|`controlPlane.hyper.port`|5050
|`controlPlane.hyper.readinessPeriodSeconds`|100
|`controlPlane.hyper.readinessTimeoutSeconds`|60
|`controlPlane.hyper.replicas`|1
|`controlPlane.hyper.svcName`|hyper
|`controlPlane.hyper.token`|token
|`controlPlane.ldap.account`|userPrincipalName
|`controlPlane.ldap.adminPassword`|-
|`controlPlane.ldap.adminUser`|-
|`controlPlane.ldap.base`|-
|`controlPlane.ldap.enabled`|false
|`controlPlane.ldap.host`|-
|`controlPlane.ldap.port`|-
|`controlPlane.ldap.ssl`|-
|`controlPlane.mpi.enabled`|true
|`controlPlane.mpi.image`|mpioperator/mpi-operator:v0.2.3
|`controlPlane.mpi.kubectlDeliveryImage`|mpioperator/kubectl-delivery:v0.2.3
|`controlPlane.mpi.registry.name`|mpi-private-registry
|`controlPlane.mpi.registry.password`|-
|`controlPlane.mpi.registry.url`|docker.io
|`controlPlane.mpi.registry.user`|-
|`controlPlane.objectStorage.cnvrgStorageAccessKey`|AKIAIOSFODNN7EXAMPLE
|`controlPlane.objectStorage.cnvrgStorageAzureAccessKey`|-
|`controlPlane.objectStorage.cnvrgStorageAzureAccountName`|-
|`controlPlane.objectStorage.cnvrgStorageAzureContainer`|-
|`controlPlane.objectStorage.cnvrgStorageBucket`|cnvrg-storage
|`controlPlane.objectStorage.cnvrgStorageEndpoint`|-
|`controlPlane.objectStorage.cnvrgStorageProject`|-
|`controlPlane.objectStorage.cnvrgStorageRegion`|eastus
|`controlPlane.objectStorage.cnvrgStorageSecretKey`|wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
|`controlPlane.objectStorage.cnvrgStorageType`|minio
|`controlPlane.objectStorage.gcpKeyfileMountPath`|/tmp/gcp_keyfile
|`controlPlane.objectStorage.gcpKeyfileName`|key.json
|`controlPlane.objectStorage.gcpStorageSecret`|gcp-storage-secret
|`controlPlane.rbac.role`|cnvrg-control-plan-role
|`controlPlane.rbac.roleBindingName`|cnvrg-control-plan-binding
|`controlPlane.rbac.serviceAccountName`|cnvrg
|`controlPlane.searchkiq.cpu`|750m
|`controlPlane.searchkiq.enabled`|true
|`controlPlane.searchkiq.killTimeout`|60
|`controlPlane.searchkiq.memory`|750Mi
|`controlPlane.searchkiq.replicas`|1
|`controlPlane.seeder.createBucketCmd`|mb.sh
|`controlPlane.seeder.image`|docker.io/cnvrg/cnvrg-boot:v0.26-tenancy
|`controlPlane.seeder.seedCmd`|rails db:migrate && rails db:seed && rails libraries:update
|`controlPlane.sidekiq.cpu`|1000m
|`controlPlane.sidekiq.enabled`|true
|`controlPlane.sidekiq.killTimeout`|60
|`controlPlane.sidekiq.memory`|3750Mi
|`controlPlane.sidekiq.replicas`|2
|`controlPlane.sidekiq.split`|true
|`controlPlane.smtp.domain`|-
|`controlPlane.smtp.password`|-
|`controlPlane.smtp.port`|-
|`controlPlane.smtp.server`|-
|`controlPlane.smtp.username`|-
|`controlPlane.systemkiq.cpu`|500m
|`controlPlane.systemkiq.enabled`|true
|`controlPlane.systemkiq.killTimeout`|60
|`controlPlane.systemkiq.memory`|500Mi
|`controlPlane.systemkiq.replicas`|1
|`controlPlane.tenancy.dedicatedNodes`|false
|`controlPlane.tenancy.enabled`|false
|`controlPlane.tenancy.key`|cnvrg-taint
|`controlPlane.tenancy.value`|true
|`controlPlane.webapp.cpu`|2000m
|`controlPlane.webapp.enabled`|true
|`controlPlane.webapp.failureThreshold`|4
|`controlPlane.webapp.image`|cnvrg/core:3.1.5
|`controlPlane.webapp.initialDelaySeconds`|10
|`controlPlane.webapp.memory`|4Gi
|`controlPlane.webapp.nodePort`|30080
|`controlPlane.webapp.oauthProxy.skipAuthRegex.0`|^\/api
|`controlPlane.webapp.oauthProxy.skipAuthRegex.1`|\/assets
|`controlPlane.webapp.oauthProxy.skipAuthRegex.2`|\/healthz
|`controlPlane.webapp.oauthProxy.skipAuthRegex.3`|\/public
|`controlPlane.webapp.oauthProxy.skipAuthRegex.4`|\/pack
|`controlPlane.webapp.oauthProxy.skipAuthRegex.5`|\/vscode.tar.gz
|`controlPlane.webapp.oauthProxy.skipAuthRegex.6`|\/gitlens.vsix
|`controlPlane.webapp.oauthProxy.skipAuthRegex.7`|\/ms-python-release.vsix
|`controlPlane.webapp.passengerMaxPoolSize`|20
|`controlPlane.webapp.port`|8080
|`controlPlane.webapp.readinessPeriodSeconds`|25
|`controlPlane.webapp.readinessTimeoutSeconds`|20
|`controlPlane.webapp.replicas`|1
|`controlPlane.webapp.svcName`|app
|`dbs.es.cpuLimit`|2000m
|`dbs.es.cpuRequest`|1000m
|`dbs.es.enabled`|true
|`dbs.es.fsGroup`|1000
|`dbs.es.image`|docker.io/cnvrg/cnvrg-es:v7.8.1
|`dbs.es.javaOpts`|-
|`dbs.es.memoryLimit`|4Gi
|`dbs.es.memoryRequest`|1Gi
|`dbs.es.nodePort`|32200
|`dbs.es.patchEsNodes`|true
|`dbs.es.port`|9200
|`dbs.es.runAsUser`|1000
|`dbs.es.serviceAccount`|es
|`dbs.es.storageClass`|-
|`dbs.es.storageSize`|30Gi
|`dbs.es.svcName`|elasticsearch
|`dbs.minio.cpuRequest`|1000m
|`dbs.minio.enabled`|true
|`dbs.minio.image`|docker.io/minio/minio:RELEASE.2020-09-17T04-49-20Z
|`dbs.minio.memoryRequest`|2Gi
|`dbs.minio.nodePort`|30090
|`dbs.minio.port`|9000
|`dbs.minio.replicas`|1
|`dbs.minio.serviceAccount`|minio
|`dbs.minio.sharedStorage.consistentHash.key`|httpQueryParameterName
|`dbs.minio.sharedStorage.consistentHash.value`|uploadId
|`dbs.minio.sharedStorage.enabled`|enabled
|`dbs.minio.storageClass`|-
|`dbs.minio.storageSize`|100Gi
|`dbs.minio.svcName`|minio
|`dbs.pg.cpuRequest`|4000m
|`dbs.pg.dbname`|cnvrg_production
|`dbs.pg.enabled`|true
|`dbs.pg.fixpg`|true
|`dbs.pg.fsGroup`|26
|`dbs.pg.hugePages.enabled`|false
|`dbs.pg.hugePages.memory`|-
|`dbs.pg.hugePages.size`|2Mi
|`dbs.pg.image`|centos/postgresql-12-centos7
|`dbs.pg.maxConnections`|100
|`dbs.pg.memoryRequest`|4Gi
|`dbs.pg.pass`|pg_pass
|`dbs.pg.port`|5432
|`dbs.pg.runAsUser`|26
|`dbs.pg.secretName`|cnvrg-pg-secret
|`dbs.pg.serviceAccount`|pg
|`dbs.pg.sharedBuffers`|64MB
|`dbs.pg.storageClass`|-
|`dbs.pg.storageSize`|80Gi
|`dbs.pg.svcName`|postgres
|`dbs.pg.user`|cnvrg
|`dbs.redis.appendonly`|yes
|`dbs.redis.enabled`|true
|`dbs.redis.image`|docker.io/cnvrg/cnvrg-redis:v3.0.5.c2
|`dbs.redis.limits.cpu`|1000m
|`dbs.redis.limits.memory`|2Gi
|`dbs.redis.port`|6379
|`dbs.redis.requests.cpu`|100m
|`dbs.redis.requests.memory`|200Mi
|`dbs.redis.serviceAccount`|redis
|`dbs.redis.storageClass`|-
|`dbs.redis.storageSize`|10Gi
|`dbs.redis.svcName`|redis
|`gpu.nvidiaDp.enabled`|true
|`gpu.nvidiaDp.image`|nvcr.io/nvidia/k8s-device-plugin:v0.9.0
|`infraNamespace`|cnvrg-infra
|`infraReconcilerCm`|infra-reconciler-cm
|`logging.elastalert.containerPort`|3030
|`logging.elastalert.cpuLimit`|400m
|`logging.elastalert.cpuRequest`|100m
|`logging.elastalert.enabled`|true
|`logging.elastalert.fsGroup`|1000
|`logging.elastalert.image`|bitsensor/elastalert:3.0.0-beta.1
|`logging.elastalert.memoryLimit`|800Mi
|`logging.elastalert.memoryRequest`|200Mi
|`logging.elastalert.nodePort`|32030
|`logging.elastalert.port`|80
|`logging.elastalert.runAsUser`|1000
|`logging.elastalert.storageClass`|-
|`logging.elastalert.storageSize`|30Gi
|`logging.elastalert.svcName`|elastalert
|`logging.enabled`|true
|`logging.fluentbit.image`|cnvrg/cnvrg-fluentbit:v1.7.2
|`logging.kibana.cpuLimit`|1000m
|`logging.kibana.cpuRequest`|100m
|`logging.kibana.enabled`|true
|`logging.kibana.image`|docker.elastic.co/kibana/kibana-oss:7.8.1
|`logging.kibana.memoryLimit`|2Gi
|`logging.kibana.memoryRequest`|100Mi
|`logging.kibana.nodePort`|30601
|`logging.kibana.port`|8080
|`logging.kibana.serviceAccount`|default
|`logging.kibana.svcName`|kibana
|`monitoring.dcgmExporter.enabled`|true
|`monitoring.dcgmExporter.image`|nvcr.io/nvidia/k8s/dcgm-exporter:2.1.4-2.3.1-ubuntu18.04
|`monitoring.enabled`|true
|`monitoring.grafana.enabled`|true
|`monitoring.grafana.image`|grafana/grafana:7.3.4
|`monitoring.grafana.nodePort`|30012
|`monitoring.grafana.oauthProxy.skipAuthRegex.0`|\/api\/health
|`monitoring.grafana.port`|8080
|`monitoring.grafana.svcName`|grafana
|`monitoring.kubeStateMetrics.enabled`|true
|`monitoring.kubeStateMetrics.image`|quay.io/coreos/kube-state-metrics:v1.9.7
|`monitoring.kubeletServiceMonitor`|true
|`monitoring.nodeExporter.enabled`|true
|`monitoring.nodeExporter.image`|quay.io/prometheus/node-exporter:v1.0.1
|`monitoring.prometheus.cpuRequest`|200m
|`monitoring.prometheus.enabled`|true
|`monitoring.prometheus.image`|quay.io/prometheus/prometheus:v2.22.1
|`monitoring.prometheus.memoryRequest`|500Mi
|`monitoring.prometheus.nodePort`|30909
|`monitoring.prometheus.port`|9090
|`monitoring.prometheus.storageClass`|-
|`monitoring.prometheus.storageSize`|50Gi
|`monitoring.prometheus.svcName`|prometheus
|`monitoring.prometheusOperator.enabled`|true
|`monitoring.prometheusOperator.images.configReloaderImage`|-
|`monitoring.prometheusOperator.images.kubeRbacProxyImage`|quay.io/brancz/kube-rbac-proxy:v0.8.0
|`monitoring.prometheusOperator.images.operatorImage`|quay.io/prometheus-operator/prometheus-operator:v0.44.1
|`monitoring.prometheusOperator.images.prometheusConfigReloaderImage`|quay.io/prometheus-operator/prometheus-config-reloader:v0.44.1
|`monitoring.upstreamPrometheus`|prometheus-operated.cnvrg-infra.svc.cluster.local:9090
|`namespaceTenancy`|false
|`networking.https.cert`|-
|`networking.https.certSecret`|-
|`networking.https.enabled`|false
|`networking.https.key`|-
|`networking.ingress.enabled`|true
|`networking.ingress.ingressType`|istio
|`networking.ingress.istioGwName`|-
|`networking.ingress.perTryTimeout`|3600s
|`networking.ingress.retriesAttempts`|5
|`networking.ingress.timeout`|18000s
|`networking.istio.enabled`|true
|`networking.istio.externalIp`|-
|`networking.istio.hub`|docker.io/istio
|`networking.istio.ingressSvcAnnotations`|-
|`networking.istio.ingressSvcExtraPorts`|-
|`networking.istio.loadBalancerSourceRanges`|-
|`networking.istio.mixerImage`|mixer
|`networking.istio.operatorImage`|docker.io/istio/operator:1.8.1
|`networking.istio.pilotImage`|pilot
|`networking.istio.proxyImage`|proxyv2
|`networking.istio.tag`|1.8.1
|`registry.name`|cnvrg-infra-registry
|`registry.password`|-
|`registry.url`|docker.io
|`registry.user`|-
|`sso.adminUser`|-
|`sso.azureTenant`|-
|`sso.clientId`|-
|`sso.clientSecret`|-
|`sso.cookieSecret`|-
|`sso.emailDomain`|-
|`sso.enabled`|false
|`sso.image`|cnvrg/cnvrg-oauth-proxy:v7.0.1.c2
|`sso.oidcIssuerUrl`|-
|`sso.provider`|-
|`sso.redisConnectionUrl`|redis://redis:6379
|`storage.enabled`|false
|`storage.hostpath.cpuLimit`|200m
|`storage.hostpath.cpuRequest`|100m
|`storage.hostpath.defaultSc`|false
|`storage.hostpath.enabled`|false
|`storage.hostpath.hostPath`|/cnvrg-hostpath-storage
|`storage.hostpath.image`|quay.io/kubevirt/hostpath-provisioner
|`storage.hostpath.memoryLimit`|200Mi
|`storage.hostpath.memoryRequest`|100Mi
|`storage.hostpath.nodeName`|-
|`storage.hostpath.reclaimPolicy`|Retain
|`storage.hostpath.storageClassName`|cnvrg-hostpath-storage
|`storage.nfs.cpuLimit`|100m
|`storage.nfs.cpuRequest`|100m
|`storage.nfs.defaultSc`|false
|`storage.nfs.enabled`|false
|`storage.nfs.image`|gcr.io/k8s-staging-sig-storage/nfs-subdir-external-provisioner:v4.0.0
|`storage.nfs.memoryLimit`|200Mi
|`storage.nfs.memoryRequest`|100Mi
|`storage.nfs.path`|-
|`storage.nfs.provisioner`|cnvrg.io/ifs
|`storage.nfs.reclaimPolicy`|Retain
|`storage.nfs.server`|-


### Build
Build docker image 
```
TAG=<docker-tag> make docker-build 
```
Push docker image
```
TAG=<docker-tag> make docker-push
```
Deploy operator
```
TAG=<docker-tag> make deploy
# use single command 
TAG=<docker-tag> make docker-build docker-push deploy
```

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
   
  

### Quick start - multiple cnvrg control plans within the same K8s cluster 

```shell
add helm deploy command here 
```
### Quick start - single cnvrg control plan

```shell
add helm deploy command here 
```

### Examples 
```shell
specs exaples goes here 
```

### Chart options 
|**key**|**default value**
| ---|---| 
|`controlPlan.baseConfig.agentCustomTag`|latest
|`controlPlan.baseConfig.checkJobExpiration`|true
|`controlPlan.baseConfig.cnvrgJobUid`|1000
|`controlPlan.baseConfig.defaultComputeConfig`|/opt/kube
|`controlPlan.baseConfig.defaultComputeName`|default
|`controlPlan.baseConfig.extractTagsFromCmd`|false
|`controlPlan.baseConfig.intercom`|true
|`controlPlan.baseConfig.passengerAppEnv`|app
|`controlPlan.baseConfig.railsEnv`|app
|`controlPlan.baseConfig.runJobsOnSelfCluster`|true
|`controlPlan.baseConfig.useStdout`|true
|`controlPlan.cnvrgRouter.enabled`|false
|`controlPlan.cnvrgRouter.image`|nginx
|`controlPlan.cnvrgRouter.nodePort`|30081
|`controlPlan.cnvrgRouter.port`|80
|`controlPlan.cnvrgRouter.svcName`|routing-service
|`controlPlan.hyper.cpuLimit`|2
|`controlPlan.hyper.cpuRequest`|100m
|`controlPlan.hyper.enableReadinessProbe`|true
|`controlPlan.hyper.enabled`|true
|`controlPlan.hyper.image`|cnvrg/hyper-server:latest
|`controlPlan.hyper.memoryLimit`|4Gi
|`controlPlan.hyper.memoryRequest`|200Mi
|`controlPlan.hyper.nodePort`|30050
|`controlPlan.hyper.port`|5050
|`controlPlan.hyper.readinessPeriodSeconds`|100
|`controlPlan.hyper.readinessTimeoutSeconds`|60
|`controlPlan.hyper.replicas`|1
|`controlPlan.hyper.svcName`|hyper
|`controlPlan.hyper.token`|token
|`controlPlan.ldap.account`|userPrincipalName
|`controlPlan.ldap.enabled`|false
|`controlPlan.objectStorage.cnvrgStorageAccessKey`|AKIAIOSFODNN7EXAMPLE
|`controlPlan.objectStorage.cnvrgStorageBucket`|cnvrg-storage
|`controlPlan.objectStorage.cnvrgStorageRegion`|eastus
|`controlPlan.objectStorage.cnvrgStorageSecretKey`|wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
|`controlPlan.objectStorage.cnvrgStorageType`|minio
|`controlPlan.objectStorage.gcpKeyfileMountPath`|/tmp/gcp_keyfile
|`controlPlan.objectStorage.gcpKeyfileName`|key.json
|`controlPlan.objectStorage.gcpStorageSecret`|gcp-storage-secret
|`controlPlan.rbac.role`|cnvrg-control-plan-role
|`controlPlan.rbac.roleBindingName`|cnvrg-control-plan-binding
|`controlPlan.rbac.serviceAccountName`|cnvrg
|`controlPlan.registry.name`|cnvrg-registry
|`controlPlan.registry.url`|docker.io
|`controlPlan.searchkiq.cpu`|750m
|`controlPlan.searchkiq.enabled`|true
|`controlPlan.searchkiq.killTimeout`|60
|`controlPlan.searchkiq.memory`|750Mi
|`controlPlan.searchkiq.replicas`|1
|`controlPlan.seeder.createBucketCmd`|mb.sh
|`controlPlan.seeder.image`|docker.io/cnvrg/cnvrg-boot:v0.26-tenancy
|`controlPlan.seeder.seedCmd`|rails db:migrate && rails db:seed && rails libraries:update
|`controlPlan.sidekiq.cpu`|1750m
|`controlPlan.sidekiq.enabled`|true
|`controlPlan.sidekiq.killTimeout`|60
|`controlPlan.sidekiq.memory`|3750Mi
|`controlPlan.sidekiq.replicas`|2
|`controlPlan.sidekiq.split`|true
|`controlPlan.systemkiq.cpu`|500m
|`controlPlan.systemkiq.enabled`|false
|`controlPlan.systemkiq.killTimeout`|60
|`controlPlan.systemkiq.memory`|500Mi
|`controlPlan.systemkiq.replicas`|1
|`controlPlan.tenancy.dedicatedNodes`|false
|`controlPlan.tenancy.enabled`|false
|`controlPlan.tenancy.key`|cnvrg-taint
|`controlPlan.tenancy.value`|true
|`controlPlan.webapp.cpu`|2
|`controlPlan.webapp.enabled`|true
|`controlPlan.webapp.failureThreshold`|4
|`controlPlan.webapp.image`|cnvrg/core:3.1.5
|`controlPlan.webapp.initialDelaySeconds`|10
|`controlPlan.webapp.memory`|4Gi
|`controlPlan.webapp.nodePort`|30080
|`controlPlan.webapp.oauthProxy.skipAuthRegex.0`|^\/api
|`controlPlan.webapp.oauthProxy.skipAuthRegex.1`|\/assets
|`controlPlan.webapp.oauthProxy.skipAuthRegex.2`|\/healthz
|`controlPlan.webapp.oauthProxy.skipAuthRegex.3`|\/public
|`controlPlan.webapp.oauthProxy.skipAuthRegex.4`|\/pack
|`controlPlan.webapp.oauthProxy.skipAuthRegex.5`|\/vscode.tar.gz
|`controlPlan.webapp.oauthProxy.skipAuthRegex.6`|\/gitlens.vsix
|`controlPlan.webapp.oauthProxy.skipAuthRegex.7`|\/ms-python-release.vsix
|`controlPlan.webapp.passengerMaxPoolSize`|20
|`controlPlan.webapp.port`|8080
|`controlPlan.webapp.readinessPeriodSeconds`|25
|`controlPlan.webapp.readinessTimeoutSeconds`|20
|`controlPlan.webapp.replicas`|1
|`controlPlan.webapp.svcName`|app
|`dbs.es.cpuLimit`|2
|`dbs.es.cpuRequest`|1
|`dbs.es.enabled`|true
|`dbs.es.fsGroup`|1000
|`dbs.es.image`|docker.io/cnvrg/cnvrg-es:v7.8.1
|`dbs.es.memoryLimit`|4Gi
|`dbs.es.memoryRequest`|1Gi
|`dbs.es.nodePort`|32200
|`dbs.es.patchEsNodes`|true
|`dbs.es.port`|9200
|`dbs.es.runAsUser`|1000
|`dbs.es.serviceAccount`|es
|`dbs.es.storageSize`|30Gi
|`dbs.es.svcName`|elasticsearch
|`dbs.minio.cpuRequest`|1
|`dbs.minio.enabled`|true
|`dbs.minio.image`|docker.io/minio/minio:RELEASE.2020-09-17T04-49-20Z
|`dbs.minio.memoryRequest`|2Gi
|`dbs.minio.nodePort`|30090
|`dbs.minio.port`|9000
|`dbs.minio.replicas`|1
|`dbs.minio.serviceAccount`|
|`dbs.minio.sharedStorage.consistentHash.key`|httpQueryParameterName
|`dbs.minio.sharedStorage.consistentHash.value`|uploadId
|`dbs.minio.sharedStorage.enabled`|enabled
|`dbs.minio.storageSize`|100Gi
|`dbs.minio.svcName`|minio
|`dbs.pg.cpuRequest`|4
|`dbs.pg.dbname`|cnvrg_production
|`dbs.pg.enabled`|true
|`dbs.pg.fixpg`|true
|`dbs.pg.fsGroup`|26
|`dbs.pg.hugePages.enabled`|false
|`dbs.pg.hugePages.size`|2Mi
|`dbs.pg.image`|centos/postgresql-12-centos7
|`dbs.pg.maxConnections`|100
|`dbs.pg.memoryRequest`|4Gi
|`dbs.pg.pass`|pg_pass
|`dbs.pg.port`|5432
|`dbs.pg.runAsUser`|26
|`dbs.pg.secretName`|cnvrg-pg-secret
|`dbs.pg.serviceAccount`|default
|`dbs.pg.sharedBuffers`|64MB
|`dbs.pg.storageSize`|80Gi
|`dbs.pg.svcName`|postgres
|`dbs.pg.user`|cnvrg
|`dbs.redis.appendonly`|yes
|`dbs.redis.enabled`|true
|`dbs.redis.image`|docker.io/cnvrg/cnvrg-redis:v3.0.5.c2
|`dbs.redis.limits.cpu`|1
|`dbs.redis.limits.memory`|2Gi
|`dbs.redis.port`|6379
|`dbs.redis.requests.cpu`|100m
|`dbs.redis.requests.memory`|200Mi
|`dbs.redis.serviceAccount`|default
|`dbs.redis.storageSize`|10Gi
|`dbs.redis.svcName`|redis
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
|`logging.elastalert.storageClass`|use-default
|`logging.elastalert.storageSize`|30Gi
|`logging.elastalert.svcName`|elastalert
|`logging.enabled`|true
|`logging.fluentbit.image`|cnvrg/cnvrg-fluentbit:v1.7.2
|`logging.kibana.cpuLimit`|1
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
|`monitoring.prometheus.cpuRequest`|1
|`monitoring.prometheus.enabled`|true
|`monitoring.prometheus.image`|quay.io/prometheus/prometheus:v2.22.1
|`monitoring.prometheus.memoryRequest`|1Gi
|`monitoring.prometheus.nodePort`|30909
|`monitoring.prometheus.port`|9090
|`monitoring.prometheus.storageClass`|use-default
|`monitoring.prometheus.storageSize`|100Gi
|`monitoring.prometheus.svcName`|prometheus
|`monitoring.prometheusOperator.enabled`|true
|`monitoring.prometheusOperator.images.kubeRbacProxyImage`|quay.io/brancz/kube-rbac-proxy:v0.8.0
|`monitoring.prometheusOperator.images.operatorImage`|quay.io/prometheus-operator/prometheus-operator:v0.44.1
|`monitoring.prometheusOperator.images.prometheusConfigReloaderImage`|quay.io/prometheus-operator/prometheus-config-reloader:v0.44.1
|`monitoring.upstreamPrometheus`|prometheus-operated.cnvrg-infra.svc.cluster.local:9090
|`networking.https.enabled`|false
|`networking.ingress.ingressType`|istio
|`networking.ingress.perTryTimeout`|3600s
|`networking.ingress.retriesAttempts`|5
|`networking.ingress.timeout`|18000s
|`networking.istio.enabled`|true
|`networking.istio.hub`|docker.io/istio
|`networking.istio.mixerImage`|mixer
|`networking.istio.operatorImage`|docker.io/istio/operator:1.8.1
|`networking.istio.pilotImage`|pilot
|`networking.istio.proxyImage`|proxyv2
|`networking.istio.tag`|1.8.1
|`registry.name`|cnvrg-registry
|`registry.url`|docker.io
|`sso.enabled`|false
|`sso.image`|cnvrg/cnvrg-oauth-proxy:v7.0.1.c2
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
|`storage.hostpath.reclaimPolicy`|Retain
|`storage.hostpath.storageClassName`|cnvrg-hostpath-storage
|`storage.nfs.cpuLimit`|100m
|`storage.nfs.cpuRequest`|100m
|`storage.nfs.defaultSc`|false
|`storage.nfs.enabled`|false
|`storage.nfs.image`|gcr.io/k8s-staging-sig-storage/nfs-subdir-external-provisioner:v4.0.0
|`storage.nfs.memoryLimit`|200Mi
|`storage.nfs.memoryRequest`|100Mi
|`storage.nfs.provisioner`|cnvrg.io/ifs
|`storage.nfs.reclaimPolicy`|Retain
|`storage.nfs.storageClassName`|cnvrg-nfs-storage

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

### `CnvrgInfra` example 
```shell
apiVersion: mlops.cnvrg.io/v1
kind: CnvrgInfra
metadata:
  name: cnvrginfra
spec:
  clusterDomain: <cluster-domain>
  registry:
    user: <user>
    password: <password>
```

### `CnvrgApp` example
```shell
apiVersion: mlops.cnvrg.io/v1
kind: CnvrgApp
metadata:
  name: cnvrgapp
  namespace: cnvrg-1
spec:
  clusterDomain: <cluster-domain>
  controlPlan:
    webapp:
      image: <cnvrg-app-image> 
    registry:
      user: <user>
      password: <password>
```
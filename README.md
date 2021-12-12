# cnvrg.io operator (v3)
---
## Deploy cnvrg stack on EKS | AKS | GKE | OpenShift | On-Premise clusters

### Architecture overview 
cnvrg operator may deploy cnvrg stack in two different ways
1. Multiple cnvrg control planes within the same cluster separated by namespaces - suitable for multi tenancy deployments  
```shell
                            ---------cnvrg infra namespace----------
                            | Cluster scope prometheus             |
                            | Prometheus node exporter             |
                            | Kube state metrics                   |
                            | Cluster scope service monitors       |     
                            | Fluentbit                            |
                            | Istio control plane                  |
                            | Storage provisioners (hostpath/nfs)  |
                            ----------------------------------------           
---------cnvrg control plane 1 namespace-------  ---------cnvrg control plane 2 namespace-------
| cnvrg control plane (webapp, sidekiqs, etc.)|  | cnvrg control plane (webapp, sidekiqs, etc.)|
| PostgreSQL                                  |  | PostgreSQL                                  |
| ElasticSearch + Kibana                      |  | ElasticSearch + Kibana                      |
| Minio                                       |  | Minio                                       |
| Redis                                       |  | Redis                                       |
| Namespace scope Prometheus + Grafana        |  | Namespace scope Prometheus + Grafana        |
| Namespace scope service monitors            |  | Namespace scope service monitors            |
| Istio Gateway + VirtualServices             |  | Istio Gateway + VirtualServices             |
-----------------------------------------------  -----------------------------------------------
                    
```
2. Single cnvrg control plane in dedicated namespace 
```shell
                        ----------------cnvrg namespace--------------------
                        | Cluster scope prometheus                        |
                        | Prometheus node exporter                        |
                        | Kube state metrics                              |
                        | Cluster scope service monitors                  |     
                        | Namespace scope service monitors                |     
                        | Fluentbit                                       |
                        | Istio control plane                             |
                        | Storage provisioners (hostpath/nfs)             |   
                        | cnvrg control plane (webapp, sidekiqs, etc.)    |
                        | PostgreSQL                                      |
                        | ElasticSearch + Kibana                          | 
                        | Minio                                           |
                        | Redis                                           |  
                        | IstioGateway + VirtualServices                  |
                        ---------------------------------------------------           
```

* [Stack requirements](./docs/requirements.md)
* [Components](./docs/components.md)
* [Quick start](./docs/quickstart.md)
* [Deployment Examples](./docs/deployments)


# Configuration

Helm chart command line options

1. [Globals](#control-plane-options)
2. [Cnvrg Control Plane options](#control-plane-options)
3. [DataBases options](#databases-options)
4. [Logging options](#logging-options)
5. [Monitoring options](#monitoring-options)
6. [Networking options](#networking-options)
7. [SSO options](#sso-options)
8. [LDAP (Active Directory)](#sso-options)
9. [Storage options](#storage-options)
10. [Tenancy options](#tenancy-options)
11. [Registry options](#registry-options)
12. [Labels and Annotations](#labels-and-annotations)
13. [Automatic Config Reload](#automatic-config-reloader)

#### Globals

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`clusterDomain` | - | DNS A wildcard record resolving to K8s Ingress IP/LoadBalancer, example: `*.cnvrg.my-org.com -> 1.2.3.4`
`spec` | `allinone` | can be set to one of `allinone` - for single namespace deployment. `infra`  and `ccp` for multi namespaces cnvrg deployments
`imageHub` | docker.io/cnvrg | the images registry

#### Control Plane options

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`controlPlane.image` |  cnvrg/core:3.6.99 | cnvrg control plane image
`controlPlane.baseConfig.agentCustomTag` |  latest | cnvrg agent image tag
`controlPlane.baseConfig.featureFlags` | {} | map of strings, usage example: `--set controlPlane.baseConfig.featureFlags.FOO="BAR"`
`controlPlane.baseConfig.intercom` |  true | set to false to disable intercom
`controlPlane.hyper.enabled` |  true | set to false to disable hyper
`controlPlane.objectStorage.type` |  minio | supported values: `minio,aws,azure,gcp`
`controlPlane.objectStorage.bucket` |  cnvrg-storage | S3 bucket name
`controlPlane.objectStorage.region` |  eastus | bucket region
`controlPlane.objectStorage.accessKey` |  - | bucket access key (if blank - auto generated)
`controlPlane.objectStorage.secretKey` |  - | bucket secret key (if blank - auto generated)
`controlPlane.objectStorage.endpoint` |  - | bucket endpoint (if blank - auto generated )
`controlPlane.objectStorage.azureAccountName` |  - | azure storage account name
`controlPlane.objectStorage.azureContainer` |  - | azure storage container name
`controlPlane.objectStorage.gcpProject` |  - | gcp project
`controlPlane.objectStorage.gcpSecretRef` |  gcp-storage-secret | gcp storage secret
`controlPlane.searchkiq.enabled` |  true | set to false to disable searchkiq
`controlPlane.sidekiq.enabled` |  true | set to false to disable sidekiq
`controlPlane.sidekiq.split` |  true | set to false to disable sidekiq split
`controlPlane.systemkiq.enabled` |  true | set to false to disable systemkiq split
`controlPlane.webapp.hap.enabled` |  true | set to false to disable hpa
`controlPlane.webapp.hap.maxReplicas` |  5 | set max replicas for HPA
`controlPlane.sidekiq.hap.enabled` |  true | set to false to disable hpa
`controlPlane.sidekiq.hap.maxReplicas` |  5 | set max replicas for HPA
`controlPlane.searchkiq.hap.enabled` |  true | set to false to disable hpa
`controlPlane.searchkiq.hap.maxReplicas` |  5 | set max replicas for HPA
`controlPlane.systemkiq.hap.enabled` |  true | set to false to disable hpa
`controlPlane.systemkiq.hap.maxReplicas` |  5 | set max replicas for HPA
`controlPlane.smtp.server` |  - | smtp server
`controlPlane.smtp.port` |  587 | smtp port
`controlPlane.smtp.username` |  - | smtp username
`controlPlane.smtp.password` |  - | smtp password
`controlPlane.smtp.domain` |  - | smtp domain
`controlPlane.smtp.opensslVerifyMode` | - | openssl verify mode for cnvrg smtp client 
`controlPlane.smtp.sender` | info@cnvrg.io | the email address of the sender   
`controlPlane.webapp.enabled` |  true | set to false to disable webapp
`controlPlane.webapp.replicas` |  1 | webapp replicas number
`controlPlane.mpi.enabled` |  true | set to false to disable mpi
`controlPlane.mpi.image` |  mpioperator/mpi-operator:v0.2.3 | mpi operator image
`controlPlane.mpi.kubectlDeliveryImage` |  mpioperator/kubectl-delivery:v0.2.3 | mpi kubectl delivery image
`controlPlane.mpi.registry.url` |  docker.io | mpi registry url
`controlPlane.mpi.registry.user` |  - | mpi registry user
`controlPlane.mpi.registry.password` | - | mpi registry password
`controlPlane.mpi.extraArgs` | {} |  map of strings, usage example: `--set controlPlane.mpi.extraArgs.FOO="BAR"`

#### DataBases options

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`dbs.es.enabled` |  true | set to false to disable elasticsearch
`dbs.es.storageSize` |  80Gi | storage size for elasticsearch
`dbs.es.storageClass` |  - | storage class, if blank default storage class will be used
`dbs.minio.enabled` |  true | set to false to disable minio
`dbs.minio.storageSize` |  100Gi | storage size for minio
`dbs.minio.storageClass` |  - | storage class, if blank default storage class will be used
`dbs.pg.enabled` |  true | set to false to disable postgresql
`dbs.pg.storageSize` |  80Gi | storage size for postgresql
`dbs.pg.storageClass` |  - | storage class, if blank default storage class will be used
`dbs.pg.hugePages.enabled` |  false | set to true to enable HubePages support for postgresql
`dbs.pg.hugePages.size` |  2Mi | size of hubePages (1Mi, 2Mi, 1Gi)
`dbs.pg.hugePages.memory` |  - | memory amount to use from the hubepages, default 4Gi
`dbs.redis.enabled` |  true | set to false to disable redis
`dbs.redis.storageSize` |  10Gi | storage size for redis
`dbs.redis.storageClass` |  - | storage class, if blank default storage class will be used

#### Logging options

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`logging.fluentbit.enabled` |  true | set to false to disable fluentbit
`logging.elastalert.enabled` |  true | set to false to disable elastalert
`logging.elastalert.storageSize` |  30Gi | storage size for elastalert
`logging.elastalert.storageClass` |  - | storage class, if blank default storage class will be used
`logging.kibana.enabled` |  true | set to false to disable kibana

#### Monitoring options

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`monitoring.dcgmExporter.enabled` |  true | set to false to disable dcgmExporter
`monitoring.nodeExporter.enabled` |  true | set to false to disable nodeExporter
`monitoring.kubeStateMetrics.enabled` |  true | set to false to disable kubeStateMetrics
`monitoring.grafana.enabled` |  true | set to false to disable grafana
`monitoring.prometheusOperator.enabled` |  true | set to false to disable prometheusOperator
`monitoring.prometheus.enabled` |  true | set to false to disable prometheus
`monitoring.prometheus.storageSize` |  50Gi | storage for Prometheus instance
`monitoring.prometheus.storageClass` |  - | storage class, if blank default storage class will be used
`monitoring.defaultServiceMonitors.enabled` |  true | set to false to disable defaultServiceMonitors
`monitoring.cnvrgIdleMetricsExporter.enabled` |  true | set to false to disable cnvrgIdleMetricsExporter

#### Networking options

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`networking.https.enabled` |  false | set to false to disable https
`networking.https.certSecret` |  - | K8s tls secret
`networking.ingress.type` |  istio | ingress type: (`istio\|ingress\|openshift\|nodeport`)
`networking.ingress.istioGwEnabled` |  true | either deploy or not Istio GW
`networking.ingress.istioGwName` |  istio-gw-[namespace] | name of the istio GW (either to use or create and use if istioGwEnabled is true)
`networking.istio.enabled` |  true | set to false to disable istio deployment
`networking.istio.externalIp` |  [] | list of IPs to use for istio ingress service: example: `--set networking.istio.externalIp={10.0.0.22,10.0.0.33}`
`networking.istio.ingressSvcExtraPorts` |  [] | list extra ports for istio ingress service: example: `--set networking.istio.externalIp={1111,2222}`
`networking.istio.lbSourceRanges` |  [] | list extra LB sources ranges, example: `--set networking.istio.externalIp={1.1.1.1/32,2.2.2.2/30}`
`networking.istio.ingressSvcAnnotations` |  {} | map of strings for Istio SVC annotations, example : `--set networking.istio.ingressSvcAnnotations=networking.istio.ingressSvcAnnotations.service\.beta\.kubernetes\.io\/aws-load-balancer-backend-protocol=tcp`
`networking.proxy.enabled` |  false | set to true when yours K8s is behind HTTP/S proxy
`networking.proxy.httpProxy` |  [] | list of http proxies to use, example `--set networking.proxy.httpProxy={http://172.17.0.5:3128}`
`networking.proxy.httpsProxy` |  [] | list of http proxies to use, example `--set networking.proxy.httpsProxy={http://172.17.0.5:3128}`
`networking.proxy.noProxy` |  .svc,.svc.cluster.local,[k8s-api-ip-calculated-automatically],127.0.0.1,kubernetes.default.svc,kubernetes.default.svc.cluster.local,localhost | list of extra no_proxy values to use (will be always appended to default list), example `--set networking.proxy.noProxy={my.extra.domain.com}`

#### SSO options

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`sso.enabled` | false |  set to true to enable sso
`sso.adminUser` |  - | cnvrg cluster admin user
`sso.provider` |  - | one of the [following](https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/oauth_provider)
`sso.emailDomain` |  [] | list of emails allowed to login
`sso.clientId` |  - | oauth2 client ID
`sso.clientSecret` |  - | oauth2 client secret
`sso.azureTenant` |  - | if `sso.provider=azure` set `azureTenant`
`sso.oidcIssuerUrl` |  - | if `sso.provider=oidc` set `oidcIssuerUrl`

#### Ldap - Active directory 

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`ldap.enabled` | false | set to true to enable sso
`ldap.host` | - | Ldap host 
`ldap.port` | - | Ldap port
`ldap.account` | - | userPrincipalName
`ldap.base` | - | for example: dc=my-domain,dc=local
`ldap.adminUser` | - | admin user
`ldap.adminPassword` | - | admin password
`ldap.ssl` | "false" | ("true" or "false")


#### Storage options

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`storage.hostpath.enabled` |  False | set to true to enable hostpath provisioner
`storage.hostpath.defaultSc` |  False | set to true to make hostpath default storage class (name: `cnvrg-hostpath-storage`)
`storage.hostpath.path` |  /cnvrg-hostpath-storage | host directory for storage
`storage.hostpath.image` |  quay.io/kubevirt/hostpath-provisioner | hostpath provisioner image
`storage.hostpath.reclaimPolicy` |  Retain | [storage class reclaim policy](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#storage-object-in-use-protection)
`storage.nfs.enabled` |  False | set to true to enable hostpath nfs client provisioner
`storage.nfs.server` |  - | Ip address of the NFS server
`storage.nfs.path` |  - | NFS export path
`storage.nfs.defaultSc` |  False | set to true to make NFS default storage class (name: `cnvrg-nfs-storage`)
`storage.nfs.image` |  gcr.io/k8s-staging-sig-storage/nfs-subdir-external-provisioner:v4.0.0 | Nfs provisioner image
`storage.nfs.reclaimPolicy` |  Retain | [storage class reclaim policy](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#storage-object-in-use-protection)

#### Tenancy options

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`tenancy.enabled` |  False | when true, ccp workloads will be scheduled only on nodes that match node selector: `purpose=cnvrg-control-plane`
`tenancy.key` |  purpose | node selector key
`tenancy.value` |  cnvrg-control-plane | node selector value

#### Registry options

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`registry.url` |  docker.io | registry for pulling images
`registry.user` |  - | registry user
`registry.password` |  - |  registry password

#### Labels and Annotations

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`lables` | `owner: cnvrg-control-plane` | key:value map of labels to be passed to every K8s resource deployed by Operator. usage example: `--set labels.foo="bar" `
`annotations` | - | key:value map of annotations to be passed to every K8s resource deployed by Operator. usage example: `--set annotations.foo="bar" `

#### Automatic Config Reloader

|**Flag**|**Default value**|**Description**
| ------------------|---|-------
`configReloader.enabled` | `true` | set to false to disable config reloader, note, once disabled, cnvrg admin has to manually restart relevant pods on configuration changes    



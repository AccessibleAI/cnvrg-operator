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
   
  

### Quick start
Setup multi tenant cnvrg cluster




### Build & Dev
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
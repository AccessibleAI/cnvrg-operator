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
* [Configuration](./docs/configuration.md)
* [Deployment Examples](./docs/deployments)




# Stack components

### Core components 
1. Cnvrg Operator
2. WebApp
3. Hyper
4. Sidekiq 
5. Searchkiq
6. Systemkiq 

#### DataBases
1. PostgreSQL
2. ElasticSearch
3. Redis
4. Minio (if `controlPlane.objectStorage.Type=minio` )

#### Logging
1. Fluentbit
2. ElastAlert
3. Kibana

#### Monitoring
1. Prometheus Operator 
2. Prometheus Instance
3. Grafana
4. Prometheus Node Exporter
5. Default Service Monitors (Node exporters, kubelet, kube-state metrics) 
6. DCGM exporter 

#### Networking 
1. Istio Operator (if `ingress.Type=istio` )
2. Istio Instance (if `ingress.Type=istio`)
3. K8s Ingress (if `ingress.Type=ingress`)
4. OpenShift Route (if `ingress.Type=openshift`)
5. NodePort (if `ingress.Type=nodeport`)

#### Storage 
1. Nfs client provisioner (if `storage.nfs.enabled=true`)
2. Hostpath provisioner (if `storage.hostpath.enabled=true`)
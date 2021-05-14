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
| PostgreSQL                                 |  | PostgreSQL                                 |
| ElasticSearch + Kibana                     |  | ElasticSearch + Kibana                     |
| Minio                                      |  | Minio                                      |
| Redis                                      |  | Redis                                      |
| Namespace scope Prometheus + Grafana       |  | Namespace scope Prometheus + Grafana       |
| Namespace scope service monitors           |  | Namespace scope service monitors           |
| Istio Gateway + VirtualServices            |  | Istio Gateway + VirtualServices            |
----------------------------------------------  ----------------------------------------------
                    
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
                        | Istio control plane                              |
                        | Storage provisioners (hostpath/nfs)             |   
                        | cnvrg control plane (webapp, sidekiqs, etc.)     |
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
  --set clusterDomain="<control-plane-domain-record>" \
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
2. Manually deploy ServiceMonitors
```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cnvrg-jobs
  namespace: openshift-monitoring
  labels:
    app: cnvrg-jobs
spec:
  jobLabel: cnvrg-job
  selector:
    matchLabels:
      exporter: cnvrg-job
  endpoints:
    - interval: 30s
      scrapeTimeout: 10s
      port: "http"
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: dcgm-exporter
  namespace: openshift-monitoring
  labels:
    app: "dcgm-exporter"
spec:
  selector:
    matchLabels:
      app: "dcgm-exporter"
  endpoints:
    - port: "metrics"
      path: "/metrics"
      interval: "15s"
``` 

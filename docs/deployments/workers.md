# Hybrid setup - Workers

#### Hybrid K8s setup cloud/on-prem - Cnvrg Control Plane on cloud and Cnvrg Workers on-prem cluster

First deploy regular CCP - this cluster will run only CCP workloads  
```bash
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
 --set clusterDomain=<cluster.domain.for.cnvrg>
``` 

Second, deploy workers cluster

```bash
helm install cnvrg cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
  --create-namespace -n cnvrg \
  --set clusterDomain=workers.<cluster.domain.for.cnvrg> \
  --set controlPlane.webapp.enabled=false \
  --set controlPlane.sidekiq.enabled=false \
  --set controlPlane.searchkiq.enabled=false \
  --set controlPlane.systemkiq.enabled=false \
  --set controlPlane.hyper.enabled=false \
  --set logging.elastalert.enabled=false \
  --set dbs.minio.enabled=false
```

Cluster domain for worker cluster have to be a subdomain of the CCP cluster.
For example: if CCP domain is `ccp.domain.com`, the worker 
cluster domain must be `workers.ccp.domain.com`
Also, DNS wildcard resolving should be configured for each domain independently.
I.e. the CCP domain `*.ccp.domain.com` should resolve to ingress IP of the CCP cluster, 
while `workers.ccp.domain.com` should resolve to ingress IP for the workers cluster. 

When running in hybrid setup, make user both sites are using same http scheme, 
either both sites should be HTTP or both should be HTTPS.

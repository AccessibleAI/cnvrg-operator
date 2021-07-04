# Multi-Tenancy deployment

#### Deploy Cnvrg Control Plane in multi-tenant cluster (CCP in each namespace)  

First deploy CnvrgInfra spec  
```bash
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg-infra --timeout 1500s \
 --set spec="infra" \
 --set clusterDomain=<infra-cluster.domain.for.cnvrg> \
 # add extra params here 
``` 

Second, deploy CCP  cluster

```bash
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n <cnvrg-namespace> --timeout 1500s \
 --set spec="ccp" \
 --set clusterDomain=<ccp-cluster.domain.for.cnvrg> \
 # add extra params here 
```


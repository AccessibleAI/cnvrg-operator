# Build-in Minio instance

#### Single K8s cluster with build-in Minio instance

```bash
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
 --set clusterDomain=<cluster.domain.for.cnvrg>
```
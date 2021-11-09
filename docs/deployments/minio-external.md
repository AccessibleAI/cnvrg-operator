# Dedicated on-prem S3

#### Single K8s cluster with dedicated on-prem S3 Object Storage (not cnvrg build-in Minio)

```bash
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
 --set clusterDomain=<cluster.domain.for.cnvrg>
 --set controlPlane.objectStorage.endpoint="<object-storage-end-point-url>" \
 --set controlPlane.objectStorage.bucket="<bucket-name>" \
 --set controlPlane.objectStorage.accessKey="<access-key>" \
 --set controlPlane.objectStorage.secretKey="<secret-key>" \
 --set dbs.minio.enabled=false
```

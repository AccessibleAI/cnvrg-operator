# Cnvrg with AWS S3 bucket
 #### Single K8s cluster with AWS S3 bucket

```bash
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
 --set clusterDomain=<cluster.domain.for.cnvrg> \
 --set controlPlane.objectStorage.type="aws" \
 --set controlPlane.objectStorage.bucket="<bucket-name>" \
 --set controlPlane.objectStorage.region="<bucket-region>" \
 --set controlPlane.objectStorage.accessKey="<access-key>" \
 --set controlPlane.objectStorage.secretKey="<secret-key>"
```
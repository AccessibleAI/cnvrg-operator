# Cnvrg with GCP Bucket

#### Single K8s cluster with single node group with GCP Bucket

```bash
# first create GCP bucket secret 
kubectl create secret generic gcp-storage-secret --from-file=key.json=/path/to/gcp/json.key -n cnvrg  

# once GCP bucket secret in place, install cnvrg
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
 --set clusterDomain=<cluster.domain.for.cnvrg> \
 --set controlPlane.objectStorage.type="gcp" \
 --set controlPlane.objectStorage.gcpProject="<gcp-project>"
```
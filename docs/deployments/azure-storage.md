# Cnvrg with Azure Storage Account
#### Single K8s cluster with Azure Storage Account 

```bash
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
 --set clusterDomain=<cluster.domain.for.cnvrg> \
 --set controlPlane.objectStorage.type="azure" \
 --set controlPlane.objectStorage.accessKey="<access-key>" \
 --set controlPlane.objectStorage.azureAccountName="<account-name>" \
 --set controlPlane.objectStorage.azureContainer="<azure-container>"
```
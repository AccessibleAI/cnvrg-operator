# Cnvrg with on-prem K8s

#### Deploy Cnvrg Control Plane on-prem K8s cluster with Istio and 2 storage provisioners.
```bash
  helm template cnvrg cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
  --create-namespace -n cnvrg  \
  --set clusterDomain=<cluster.domain.for.cnvrg> \
  --set networking.istio.externalIp="{10.0.0.22}" \
  --set dbs.es.storageClass="cnvrg-hostpath-storage" \
  --set 'dbs.es.nodeSelector.kubernetes\.io/hostname=node1' \
  --set dbs.minio.storageClass="cnvrg-hostpath-storage" \
  --set 'dbs.minio.nodeSelector.kubernetes\.io/hostname=node1' \
  --set dbs.pg.storageClass="cnvrg-hostpath-storage" \
  --set 'dbs.pg.nodeSelector.kubernetes\.io/hostname=node1' \
  --set dbs.redis.storageClass="cnvrg-hostpath-storage" \
  --set 'dbs.redis.nodeSelector.kubernetes\.io/hostname=node1' \
  --set logging.elastalert.storageClass="cnvrg-hostpath-storage" \
  --set 'logging.elastalert.nodeSelector.kubernetes\.io/hostname=node1' \
  --set monitoring.prometheus.storageClass="cnvrg-hostpath-storage" \
  --set 'monitoring.prometheus.nodeSelector.kubernetes\.io/hostname=node1' \
  --set storage.hostpath.enabled=true \
  --set storage.hostpath.defaultSc=false \
  --set storage.hostpath.path="/cnvrg-hostpath-storage" \
  --set storage.nfs.enabled=true \
  --set storage.nfs.defaultSc=true \
  --set storage.nfs.server="node0" \
  --set storage.nfs.path="/mnt/nfs"
```
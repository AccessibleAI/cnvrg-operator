# Existing Istio instance with Cnvrg

#### Using existing Istio instance with Cnvrg Control Plane
```bash
helm install cnvrg cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
  --create-namespace -n cnvrg  \
  --set clusterDomain=<cluster.domain.for.cnvrg> \
  --set networking.ingress.type=istio \
  --set networking.istio.enabled=false
```
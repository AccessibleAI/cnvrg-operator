# Existing K8s Ingress controller with Cnvrg

#### Using existing K8s Ingress controller (Nginx ingress, etc...) with Cnvrg Control Plane 
```bash
helm install cnvrg cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
  --create-namespace -n cnvrg  \
  --set clusterDomain=<cluster.domain.for.cnvrg> \
  --set networking.ingress.type=ingress \
  --set networking.istio.enabled=false
```

# Cnvrg with NodePort  

#### Using K8s NodePort with Cnvrg Control Plane (for demos or small devs envs)

```bash
helm install cnvrg cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
  --create-namespace -n cnvrg  \
  --set clusterDomain=<cluster.domain.for.cnvrg> \
  --set networking.ingress.type=nodeport \
  --set networking.istio.enabled=false
```

Nodeport not suitable for production cluster. The node port should be used only for devs/poc clusters. 

The `<node-ip>` value will depend on your cluster and network setup. 
For on-prem clusters, `node-ip` usually will be `INTERNAL-IP` of one of yours K8s nodes.   
For cloud based setups, `node-ip` might be external IP of node VM or even IP of ingress service. 
To use node port, you've to have flat IP network access to your K8s cluster 
or, if you behind firewall, router, VPN or nat gateway, you've to make sure all 
[the K8s node port ranges](https://kubernetes.io/docs/concepts/services-networking/service/#nodeport) 
are open between your local machine and remote K8s cluster (same is true even for simple Minikube/Kind/K3s devs clusters)   

Please consult with your network administrator before deploying CCP with node port.

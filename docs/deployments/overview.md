# Overview
Cnvrg control plane can be deployment in different ways, 
and by any K8s native tolling. 
Our goal, is to allow you to build the best suitable AI/ML platform for your needs.  
Many stack components can be disabled or integrated into existing solutions. 

The following components might be disabled/enabled 
1. Use your own S3 Object storage  
2. Use your own Prometheus Operator or Prometheus Instance
3. Cnvrg might connect to external PostgreSQL instance 
4. Use your own Istio or other K8s compatible Ingress controller 
5. Use cnvrg's Hostpath or Nfs client provisioners

The stack deployment and management can be done by 
1. [Helm](https://helm.sh/docs/intro/install/) 
2. [GitOps engines](https://argoproj.github.io/argo-cd/) ( mainly suitable for complex deployments, based on vanilla K8s manifests)
3. [cnvrgctl](https://github.com/accessibleAI/cnvrgctl) - coming soon 

Before proceeding to actual Cnvrg Control Plane deployment, make sure the K8s clusters aligned with all the [requirements](/requirements.md).

Cnvrg allow you to compose different clusters layouts, for different needs.
1. Single K8s cluster with single node group
2. Single K8s cluster with different nodes groups  
3. Worker clusters - K8s clusters for hybrid environments (cloud and on-prem) 
4. Fully multi tenant K8s clusters with multiple Cnvrg Control Planes in each namespace

From the above options, can be constructed more complex architectures, for example, 
* Single K8s cluster with single node group + many workers cluster 
* Full multi tenant k8s cluster with worker cluster for each tenant

To compose the best fit cluster, consult our [sales teams](https://cnvrg.io/demo/) 
or try and deploy cnvrg control plane in different layouts, check [configuration section](/configuration.md).   

#### Add cnvrgv3 helm charts repo
```bash
helm repo add cnvrgv3 https://charts.v3.cnvrg.io
helm repo update 
helm search repo cnvrgv3/cnvrg -l
```
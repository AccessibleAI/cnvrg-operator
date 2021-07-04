# Stack requirements

Cnvrg Control Plane (CCP) designed and developed as a Kubernetes (K8s) native stack.
You can deploy CCP in any place where K8s can run, in cloud K8s (AWS, GCP, Azure) or 
on-premise K8s (Kubeadm,Rancher, Kubespray, etc...). K8s cluster for CCP must
include a **storage provisioner** (Block storage - **preferred**  of Filesystem storage) 
and must include an **ingress controller**.        

### K8s distro
Cloud based distros 
* AWS - EKS 
* GCP - GKE
* Azure - AKS

On-premises distros 
* Kubeadm
* Rancher (RKE v1)
* OpenShift (v4 +)
* Kubespray
* DeepOps
* Minikube

### K8s storage provisioner
Any K8s SCI compatible storage [provisioner](https://kubernetes.io/docs/concepts/storage/storage-classes/#provisioner).



CCP K8s Operator might deploy out-of-the-box two storage provisioners (**mainly useful for on-prem deployments**)
* [Hostpath SCI provisioner](https://github.com/kubevirt/hostpath-provisioner)
* [NFS client SCI provisioner](https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner)

For more details click [here](/docs/deployments/storage.md) 


### K8s ingress controller and DNS wildcard record 
Any [K8s ingress controller](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/)

DNS wildcard record which will resolve the ingress IP to the CCP cluster domain. e.g `*.cnvrg.my-org.com -> 1.2.3.4`


CCP K8s Operator might deploy out-of-the-box 
[Istio](https://istio.io/) and configure [Istio Ingress Gateway](https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/) to allow ingress traffic into the cluster.

If you don't have any DNS name for the CCP deployment, you might use [nip.io](https://nip.io) 
(**mainly useful for PoC or dev environments**)


### Minimal default compute/memory/storage requirements for production CCP deployment
|**Workload**|**Replicas**|**CPU Request**|**Memory Request**|**Storage**
| ------------------|---|-------|-----------|-----
|`webapp`           | 1 | 2000m | 4Gi       | -  
|`sidekiq`          | 2 | 2000m | 8Gi       | - 
|`searchkiq`        | 1 | 750m  | 750Mi     | - 
|`systemkiq`        | 1 | 500m  | 500Mi     | - 
|`hyper`            | 1 | 100m  | 200Mi     | -  
|`postgres`         | 1 | 4000m | 4Gi       | 80Gi 
|`redis`            | 1 | 100m  | 200Mi     | 10Gi 
|`elasticsearch`    | 1 | 100m  | 1Gi       | 80Gi 
|`kibana`           | 1 | 100m  | 200Mi     | -
|`elastalert`       | 1 | 100m  | 200Mi     | 30Gi
|`grafana`          | 1 | 100m  | 200Mi     | -
|`prometheus`       | 1 | 200m  | 500Mi     | 50Gi
|`sys-compenents`   | n | 2000m| 4Gi     | -
|**TOTAL**| **15** | **~14** | **~22Gi** | ~250Gi | - | - 


Please note, this is the minimal defaults, e.i K8s resources requests. 
The real resources values will be highly depends on the actual load and usage of CCP

Also, please note these compute resources are for CCP only. 
The actual compute power for ML workloads have to be calculated independently 
and include totals of CPU/GPU, Memory and storage expected to be consumed by the ML workloads.    




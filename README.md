# cnvrg.io operator (v5)
---

## Deploy cnvrg stack on EKS | AKS | GKE | OpenShift | On-Premise clusters

## Add cnvrg helm chart repo

```bash
helm repo add cnvrg https://charts.slim.cnvrg.io/cnvrg
```

## Deploy cnvrg MLOps

```bash
# simple deploy  
helm install cnvrg cnvrg/mlops \
 --create-namespace -n cnvrg \
 --set operatorVersion="<OPERATOR-VERSION>" \
 --set clusterDomain="<CLUSTER-DOMAIN>" \
 --set controlPlane.image="<MLOPS-APP-IMAGE>" \
 --set registry.user="<CNVRG-USERNAME>" \
 --set registry.password="<CNVRG-PASSWORD>" \
 --set controlPlane.baseConfig.agentCustomTag="<AGENT-CUSTOM-TAG>"
```
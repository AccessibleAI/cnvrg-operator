# Quickstart

### Install Cnvrg Control Plane (CCP) with Helm 

#### Add cnvrgv3 helm charts repo
```bash
helm repo add cnvrgv3 https://charts.v3.cnvrg.io
helm repo update 
helm search repo cnvrgv3/cnvrg -l
```


#### Default Cnvrg Control Plane deployment
```bash
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
 --set clusterDomain=<cluster.domain.for.cnvrg>
```


The default deployment command assume yours K8s cluster already 
have working K8s storage provisioner and Ingress controller.
This is default setup for most cloud based K8s distros (EKS/GKE/AKS).
If you are running on-prem K8s clusters, consult with your cluster administrator about the Storage provisioner and Ingress configuration.
For dev clusters, (e.i Minikube), you can enable [Ingress controller](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/).
Also, make sure the [Dynamic storage provisioner](https://minikube.sigs.k8s.io/docs/handbook/persistent_volumes/#dynamic-provisioning-and-csi) is enabled on your Minikube cluster.


**<cluster.domain.for.cnvrg>** is a wildcard DNS record which is resolving to the cluster's Ingress IP.
The Ingress IP is implemented in different ways in each compute cloud and for each K8s distro.
For example, in AWS - EKS, Ingress IP will be set to CNAME of EC2 ELB. 
In Azure AKS, real public IP will be allocated.
In on-premises K8s, cluster administrator we'll have to allocated external IP manually, or by installing [third party 
software](https://metallb.universe.tf/).

For further information about configuring Ingress traffic, consult with your cluster administrator.


### More examples

#### Deploy cnvrg with enable HTTPS (cnvrg's Istio instance TLS offload)
```bash
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
 --set clusterDomain=<cluster.domain.for.cnvrg> \
 --set networking.https.enabled=true \
 --set networking.https.certSecret="<k8s tls secret name>"
```

#### Deploy cnvrg with SSO 
```bash
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
 --set clusterDomain=<cluster.domain.for.cnvrg> \
 --set registry.user="<cnvrg-licence-user>" \
 --set registry.password="cnvrg-licence-password>" \
 --set sso.enabled=true \
 --set sso.adminUser="<cnvrg admin user>" \
 --set sso.provider="oidc" \
 --set sso.emailDomain="{<your email domain>}" \
 --set sso.clientId="<client-id>" \
 --set sso.clientSecret="<client-secret>" \
 --set sso.oidcIssuerUrl="<optional-issuer-url>" \  
 --set sso.azureTenant="<optional-azure-tenant>"  
```

#### cnvrg SSO support out-of-the-box the following list of IDP
```shell
Google Auth Provider
Azure Auth Provider
Facebook Auth Provider
GitHub Auth Provider
Keycloak Auth Provider
GitLab Auth Provider
LinkedIn Auth Provider
Microsoft Azure AD Provider
OpenID Connect Provider
login.gov Provider
Nextcloud Provider
DigitalOcean Auth Provider
Bitbucket Auth Provider
Gitea Auth Provider
```

#### Deploy cnvrg with Nfs provisioner 
```bash
helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
 --set clusterDomain=<cluster.domain.for.cnvrg> \
 --set storage.nfs.enabled=true \
 --set storage.nfs.defaultSc=true \
 --set storage.nfs.server="nfs-server" \
 --set storage.nfs.path="nfs-export"
```
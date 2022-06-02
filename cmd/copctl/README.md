# Cnvrg Operator CTL - COPCTl

### Motivation 
Cnvrg is a fullstack AI/MLOps platform. 
The fullstack in our context meaning, in addition to the **cnvrg control plane** 
we'll provide out-of-the-box the following components 
* K8s cluster wide monitoring based on [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus) project
* K8s cluster wide logging based on fluentbit & elasticsearch 
* K8s cluster wide networking based on Istio 

Unfortunately (or luckily) the fullstack approach can't be accepted when customers 
have its own monitoring, logging & networking stacks.  

To address this problem, we created the `copctl` tool, which at the end, should help 
you to deploy the cnvrg control plane without the cnvrg infra components. 
In addition, the `copctl` will allow you to configure & manage the cnvrg control plane 
in a way that perfectly match your environment limitation, with `copctl` you've all the deployment configuration 
granularity you might need.


### How it's works 
With `copctl` you can dump all the cnvrg control plane K8s manifests into a local folder. 
Once dump is done, you can start and edit the manifests in a way you like.
For example, you can change our default permissions, you can change the amount of replicas,
you can configure cnvrg control plane to use external PG, ElasticSearch, and many more.

### Quick start 
Get the `copctl` binary, currently we've a pre-build binaries for Linux and Mac (ARM).

To export the manifest, you'll need first prepare the following:
1. Wildcard domain that will be used by cnvrg control plane
2. The cnvrg control plane image 
3. The cnvrg docker registry username and password to pull the image
4. Specify the K8s ingress type, we support 3 ingress types: `istio|ingress|openshift`
5. Container runtime, we support 3 cris: `docker|containerd|cri-o`

Once you've all the details, run the following command to dump the manifests:

```shell
copctl profile dump control-plane \
  --wildcard-domain <wildcard-domain> \
  --control-plane-image <cnvrg-image> \
  --registry-user <registry-username> \
  --registry-password <registry-password> \
  --ingress <ingress-type> \
  --cri <cri-type>
```
To see all available options run `copctl profile dump control-plane -h`

**Important note**
1. You must create the K8s namespace manually ahead applying the manifests
2. You can't apply the manifests from two different dump command.
For example, you do a dump, then apply the manifests, then you are changing one of the 
dump parameters (for example, the image), then you can't apply the manifests again.
You must completely purge all the deployed manifests, and then apply a new generated manifests again. 
3. You must ensure that `cnvrg-db-init` config map does not exist in your namespace
before you applying the manifests. If you'll apply the manifests and teh `cnvrg-db-init`
config map exists, you'll get broken environment. So, make sure to delete this configmap
each time when you are deploying a fresh dumped manifests.

**Manifests directory**

Manifests directory structure by default will be flat. 
Keeping all the manifests in the same flat directory allowing you easily 
deploy the manifests on the K8s cluster by running `kubectl apply -f ./cnvrg-manfiests/` 
command.
However, if you'd like to apply the manifests steps by steps, or if you'll 
need differentiate the manifests by groups, you might want to set 
the `--preserve-templates-dir=true`. This flag will keep the original template
directory.

For example, by default, all the ~70 manifests will be dumped into `./cnvrg-manifestrs` 
directory.
If you'll set the `--preserve-templates-dir=true`, the manifests will be grouped 
by directories, so you should get something similar to this output. 
```shell
./cnvrg-manifests/
└── tmpl
    ├── apps-class.yaml
    ├── conf
    │   ├── cm
    │   │   ├── config-base.yaml
    │   │   ├── config-labels.yaml
    │   │   ├── config-networking.yaml
    │   │   ├── secret-base.yaml
    │   │   ├── secret-ldap.yaml
    │   │   ├── secret-object-storage.yaml
    │   │   └── secret-smtp.yaml
    │   └── rbac
    │       ├── buildimage-job-role.yaml
    │       ├── buildimage-job-rolebinding.yaml
    │       ├── buildimage-job-sa.yaml
    │       ├── ccp-role.yaml
    │       ├── ccp-rolebinding.yaml
    │       ├── ccp-sa.yaml
    │       ├── job-role.yaml
    │       ├── job-rolebinding.yaml
    │       ├── job-sa.yaml
    │       ├── privileged-job-role.yaml
    │       ├── privileged-job-rolebinding.yaml
    │       └── spark-job-sa.yaml
    ├── es
    │   ├── cm.yaml
    │   ├── role.yaml
    │   ├── rolebinding.yaml
    │   ├── sa.yaml
    │   ├── secret.yaml
    │   ├── sts.yaml
    │   ├── svc.yaml
    │   └── vs.yaml
    ├── hyper
    │   ├── dep.yaml
    │   └── svc.yaml
    ├── ingress
    │   └── gw.yaml
    ├── jobs-class.yaml
    ├── minio
    │   ├── dep.yaml
    │   ├── pdb.yaml
    │   ├── pvc.yaml
    │   ├── role.yaml
    │   ├── rolebinding.yaml
    │   ├── sa.yaml
    │   ├── svc.yaml
    │   └── vs.yaml
    ├── pg
    │   ├── dep.yaml
    │   ├── pdb.yaml
    │   ├── pvc.yaml
    │   ├── role.yaml
    │   ├── rolebinding.yaml
    │   ├── sa.yaml
    │   ├── secret.yaml
    │   └── svc.yaml
    ├── prometheus
    │   └── instance
    │       └── credsec.yaml
    ├── redis
    │   ├── dep.yaml
    │   ├── pdb.yaml
    │   ├── pvc.yaml
    │   ├── role.yaml
    │   ├── rolebinding.yaml
    │   ├── sa.yaml
    │   ├── secret.yaml
    │   └── svc.yaml
    ├── secret.yaml
    ├── sidekiqs
    │   ├── searchkiq-pdb.yaml
    │   ├── searchkiq.yaml
    │   ├── sidekiq-pdb.yaml
    │   ├── sidekiq.yaml
    │   ├── systemkiq-pdb.yaml
    │   └── systemkiq.yaml
    └── webapp
        ├── dep.yaml
        ├── oauth.yaml
        ├── oauthtoken.yaml
        ├── pdb.yaml
        ├── svc.yaml
        └── vs.yaml

14 directories, 70 files
```









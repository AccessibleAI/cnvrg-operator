# cnvrg.io operator (v5)
---

## Deploy cnvrg stack on EKS | AKS | GKE | OpenShift | On-Premise clusters

# Install instructions

<aside>
⚠️ **for non OCP deployments:**
first deploy istio ingress controller (using the instructions we have in notion)

then set the wildcard domain + TLS cert if needed.

</aside>

**add helm repos:**

```bash
helm repo add cnvrg-operator https://charts.slim.cnvrg.io/cnvrg-operator;
helm repo add cnvrg-cap https://charts.slim.cnvrg.io/cnvrg-cap
```

**update if already added:**

```bash
helm repo update cnvrg-operator cnvrg-cap
```

**install 1st helm repo, cnvrg-operator:**

<aside>
ℹ️ this chart will be deployed to `cnvrg-system` namespace,
and it hosts the cnvrg operator pod only.

</aside>

```bash
helm install cnvrg-operator cnvrg-operator/cnvrg-operator \
-n cnvrg-system --create-namespace \
  --set registry.user=cnvrghelm \
  --set registry.password=cabbecc7-4330-47b6-85c6-ea0ad5019cfa \
  --debug
```

**install 2nd helm repo, cnvrg-cap (the actual cnvrg instance)**

<aside>
ℹ️ this chart will be deployed to the `cnvrg` namespace,
and it hosts the rest of the cnvrg ctrl plain.

</aside>

set var for latest app image:

```bash
CNVRG_APP_IMAGE="cnvrg/app:v4.7.108-DEV-15824-cnvrg-agnostic-infra-merge-with-master-2"
```

for openshift:

```ruby
helm install cnvrg cnvrg-cap/cnvrg-cap \
  -n cnvrg --create-namespace \
  --set clusterDomain=slim4mj.gcpops.cnvrg.io \
  --set controlPlane.image=$CNVRG_APP_IMAGE \
  --set controlPlane.baseConfig.agentCustomTag=agnostic-logs \
  --set registry.user=cnvrghelm \
  --set registry.password=cabbecc7-4330-47b6-85c6-ea0ad5019cfa \
  --set networking.https.enabled="true" \
  --set networking.ingress.type=openshift \
  --set controlPlane.baseConfig.featureFlags.OCP_ENABLED="true" \
  --debug
```

for istio ingress controller:

```bash
helm install cnvrg cnvrg-cap/cnvrg-cap \
  -n cnvrg --create-namespace \
  --set clusterDomain=slim4mj.gcpops.cnvrg.io \
  --set controlPlane.image=$CNVRG_APP_IMAGE \
  --set controlPlane.baseConfig.agentCustomTag=agnostic-logs \
  --set registry.user=cnvrghelm \
  --set registry.password=cabbecc7-4330-47b6-85c6-ea0ad5019cfa \
  --set networking.ingress.istioIngressSelectorValue=ingress \
  --debug
```

for nginx ingress:

```bash
helm install cnvrg cnvrg-cap/cnvrg-cap \
  -n cnvrg --create-namespace \
  --set clusterDomain=slim4mj.gcpops.cnvrg.io \
  --set controlPlane.image=$CNVRG_APP_IMAGE \
  --set controlPlane.baseConfig.agentCustomTag=agnostic-logs \
  --set registry.user=cnvrghelm \
  --set registry.password=cabbecc7-4330-47b6-85c6-ea0ad5019cfa \
  --set networking.https.enabled="true" \
  --set networking.https.certSecret=dns-domain-tls-secret \
  --set networking.ingress.type=ingress \
  --debug
```

### SSO:

**OIDC (Keycloak)**

```bash
  --set sso.enabled=true \
  --set sso.jwks.enabled=true \
  --set sso.proxy.enabled=true \
  --set sso.pki.enabled=true \
  --set sso.central.enabled=true \
  --set sso.central.adminUser=$OIDC_ADMIN_USER \
  --set sso.central.provider=oidc \
  --set sso.central.clientId=$OIDC_CLIENT_ID \
  --set sso.central.clientSecret=$OIDC_CLIENT_SECRET \
  --set sso.central.oidcIssuerUrl=http://$KEYCLOAK_URL/realms/$KEYCLOAK_REALM
```

**Azure AD**

```bash
sso:
    central:
      emailDomain:
        - cnvrg.io
      enabled: true
      clientId: $AZURE_CLIENT_ID
      provider: azure
      adminUser: $AZURE_ADMIN_USER
      clientSecret: $AZURE_CLIENT_SECRET
      oidcIssuerUrl: https://login.microsoftonline.com/$AZURE_TENANT_ID/v2.0
```

**restart ctrl plain pods after changing SSO configs in cnvrgapp:**

```bash
kc delete cm proxy-config;
kc delete secret cp-base-secret;
kc rollout restart deploy app sidekiq searchkiq systemkiq cnvrg-jwks cnvrg-proxy-central sso-central
```

[DEV 19600 secure routes installation](https://www.notion.so/DEV-19600-secure-routes-installation-9f5702cf445c4a52871c5bebf003640f?pvs=21)

# Major changes

- doesn’t deploy istio ingress controller by default
- doesn’t deploy `kube-state-metrics` (i.e no `kubectl top pod` metrics)
- doesn’t deploy following daemonSets:
nvidia-device-plugin
dcgm-exporter
node-exporter (how metrics are exported then?)
- new prometheus version 2.37, ~~need to adjust queries in app (from 2.22)~~

# config reloader deprecated

in order to see which components are updated dynamically according to changes made in the cnvrgapp CR, check the following:

![Untitled](slim%20operator%206cd04e25d8e34a448186143912419a90/Untitled.png)

check this table for reference:

|  | owned by cnvrgapp CR | not owned by cnvrgapp CR |
| --- | --- | --- |
| “updatable” annotation = true | component will be changed dynamically once any related changes applied in cnvrgapp | unrelated to operator |
| “updatable” annotation = false | component needs to be deleted, then operator will re-create it with the values in cnvrgapp | unrelated to operator |

# Airgap Installs

### Feature Flags to set

Dependency installation will only use the custom pypi server set in organization settings.

```jsx
CNVRG_ENABLE_PIP_EXTRA_INDEX: "false"
```

Implication: disables loading of external JS file, which slows down the loading of all pages

```jsx
AIR_GAPPED: "true"
```

Ensure that intercom is disabled during install. This can be completed by adding the following to the values file and adding a feature flag listed below

```jsx
controlPlane:
  baseConfig:
    intercom: “false"
```

```jsx
SHOW_INTERCOM: "false"
```

### Example of airgapped values file

You need to define the customers docker registry for the cap install, app image and `imageHub`

```jsx
clusterDomain: aws.dilerous.cloud
controlPlane:
  image: harbor.dilerous.cloud/cnvrg/app:v4.7.52-DEV-15824-cnvrg-agnostic-infra-42-develop
  baseConfig:
    featureFlags:
      CNVRG_ENABLE_PIP_EXTRA_INDEX: "false"
      AIR_GAPPED: "true"
      SHOW_INTERCOM: "false"
      OCP_ENABLED: "true"
imageHub: harbor.dilerous.cloud/cnvrg
networking:
  ingress:
    type: openshift
registry:
  name: cnvrg-app-registry
  password: xxxxxxxxx
  url: harbor.dilerous.cloud/cnvrg
  user: admin
```

### Upload all the docker images to the customers docker registry

[](https://github.com/AccessibleAI/cnvrg_delivery/tree/main/scripts/export_import_images)

### Update cnvrg registry to point to customers internal docker registry

### Update SETTINGS to reference the customers pypi server

![Untitled](slim%20operator%206cd04e25d8e34a448186143912419a90/Untitled%201.png)

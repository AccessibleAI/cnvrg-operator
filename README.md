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

## Using external secret for SMTP server
It's an option to specify external secret for SMTP server credintials instead setting it in helm chart values or cnvrgapp CRD . 
The parameter to reference the secret is `controlPlane.smtp.CredentialsSecretRef` and the keys in the secret should be `username` and `password`.

```bash
helm install cnvrg cnvrg/mlops \
 --create-namespace -n cnvrg \
 --set controlPlane.smtp.credentialsSecretRef="SECRET-NAME"
```
secret example
```bash
apiVersion: v1
kind: Secret
metadata:
  name: SECRET-NAME
  namespace: cnvrg
type: Opaque
data:
  username: YWRtaW4=
  password: c2VjcmV0
```

## Using external secret for OAuth2 client configuration

It's an option to specify external secret for OAuth2 client configuration instead setting it in helm chart values or cnvrgapp CRD. The parameter to reference the secret is `sso.central.credentialsSecretRef` and the keys in the secret should be `clientId`, `clientSecret`

```bash
helm install cnvrg cnvrg/mlops \
 --create-namespace -n cnvrg \
 --set sso.central.credentialsSecretRef="SECRET-NAME"
```

secret example
```bash
apiVersion: v1
kind: Secret
metadata:
  name: SECRET-NAME
  namespace: cnvrg
type: Opaque
data:
  clientId: YWRtaW4=
  clientSecret: c2VjcmV0
```

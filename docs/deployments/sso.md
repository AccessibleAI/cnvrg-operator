# Cnvrg with SSO (with OAuth2)

#### Deploy Cnvrg Control Plane with SSO (with OAuth2)
Cnvrg Control Plane supports out-of-the-box following OAuth2 providers
* Google Auth Provider
* Azure Auth Provider
* Facebook Auth Provider
* GitHub Auth Provider
* Keycloak Auth Provider
* GitLab Auth Provider
* LinkedIn Auth Provider
* Microsoft Azure AD Provider
* OpenID Connect Provider
* login.gov Provider
* Nextcloud Provider
* DigitalOcean Auth Provider
* Bitbucket Auth Provider
* Gitea Auth Provider

#### Example: SSO with Keycloak

```bash
helm install cnvrg cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
  --create-namespace -n cnvrg  \
  --set clusterDomain=<cluster.domain.for.cnvrg> \
  --set sso.enabled=true \
  --set sso.adminUser="<cnvrg-admin-user>@<your-email-domain>" \
  --set sso.provider="oidc" \
  --set sso.emailDomain="{<your-email-domain>}" \
  --set sso.clientId="<client-id>" \
  --set sso.clientSecret="<client-secret>" \
  --set sso.oidcIssuerUrl="<keycloak-base-url>/auth/realms/<keycloak-realm>"
```

#### Example: SSO with AzureAD

```bash
helm install cnvrg cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s \
  --create-namespace -n cnvrg  \
  --set clusterDomain=<cluster.domain.for.cnvrg> \
  --set sso.enabled=true \
  --set sso.adminUser="<cnvrg-admin-user>@<your-email-domain>" \
  --set sso.provider="azure" \
  --set sso.emailDomain="{<your-email-domain>}" \
  --set sso.clientId="<client-id>" \
  --set sso.clientSecret="<client-secret>" \
  --set sso.azureTenant="<azure-tenant-id>"
```

Some Auth Providers require explicitly set redirect URIs for each service (e.g. Azure), 
while others allow you to use wildcard (e.g. Keycloak).
If your Auth Provider doesn't allow `*` as redirect URIs, please configure the following 
1. Redirect URI for WebApp: `http<s>://app.<cluster.domain.for.cnvrg>/oauth2/callback`
2. Redirect URI for Kibana: `http<s>://kibana.<cluster.domain.for.cnvrg>/oauth2/callback`
3. Redirect URI for Grafana: `http<s>://grafana.<cluster.domain.for.cnvrg>/oauth2/callback`

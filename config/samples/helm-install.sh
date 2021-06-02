
helm install cnvrg /tmp/cnvrg-3.1.0-dirty-76fd315.tgz -n cnvrg-infra --create-namespace --timeout 1500s \
  --set spec="infra" \
  --set clusterDomain="tenancy-infra.azops.cnvrg.io" \
  --set registry.user="cnvrghelm" \
  --set registry.password="0f827056-c363-41ad-b58b-1e40bf1f2ab2" \
  --set networking.https.enabled=true \
  --set networking.https.certSecret="cnvrg-infra" \
  --set networking.istio.externalIp="{10.0.0.22}" \
  --set sso.enabled=true \
  --set sso.adminUser="dima@cnvrg.io" \
  --set sso.provider="oidc" \
  --set sso.emailDomain="{cnvrg.io}" \
  --set sso.clientId="tenancy-demo-client" \
  --set sso.clientSecret="f918ab57-d40e-412a-ae4d-a29c79b2a43f" \
  --set sso.oidcIssuerUrl="http://keycloak.cnvrg.io/auth/realms/tenancy-demo" \
  --set storage.nfs.enabled=true \
  --set storage.nfs.defaultSc=true \
  --set storage.nfs.server="node0" \
  --set storage.nfs.path="/mnt/nfs"

helm install cnvrg /tmp/cnvrg-3.1.0-dirty-76fd315.tgz -n cnvrg-1 --create-namespace --timeout 1500s \
  --set spec="ccp" \
  --set clusterDomain="tenancy-cnvrg-1.azops.cnvrg.io" \
  --set controlPlane.image="app:master-7220" \
  --set controlPlane.baseConfig.featureFlags.MULTI_TENANT="true" \
  --set networking.https.enabled=true \
  --set networking.https.certSecret="cnvrg-1" \
  --set registry.user="cnvrghelm" \
  --set registry.password="0f827056-c363-41ad-b58b-1e40bf1f2ab2" \
  --set sso.enabled=true \
  --set sso.adminUser="dima@cnvrg.io" \
  --set sso.provider="oidc" \
  --set sso.emailDomain="{cnvrg.io}" \
  --set sso.clientId="tenancy-demo-client" \
  --set sso.clientSecret="f918ab57-d40e-412a-ae4d-a29c79b2a43f" \
  --set sso.oidcIssuerUrl="http://keycloak.cnvrg.io/auth/realms/tenancy-demo"


helm install cnvrg cnvrgv3/cnvrg -n cnvrg --create-namespace --timeout 1500s \
  --set clusterDomain="tenancy-cnvrg-1.azops.cnvrg.io" \
  --set controlPlane.image="app:master-7457" \
  --set registry.user="cnvrghelm" \
  --set registry.password="0f827056-c363-41ad-b58b-1e40bf1f2ab2" \
  --set networking.istio.externalIp="{10.0.0.22}" \
  --set storage.nfs.enabled=true \
  --set storage.nfs.defaultSc=true \
  --set storage.nfs.server="node0" \
  --set storage.nfs.path="/mnt/nfs"

kind: Secret
apiVersion: v1
metadata:
  name: proxy-config
  namespace: {{ .Namespace }}
  annotations:
    mlops.cnvrg.io/default-loader: "false"
    mlops.cnvrg.io/own: "true"
    mlops.cnvrg.io/updatable: "false"
stringData:
  conf: |-
    http_address = "0.0.0.0:8080"
    upstreams = [
        "http://127.0.0.1:8000",
         "file:///cnvrg-static/#/opstatic/"
    ]
    skip_auth_regex = [
        "^\/opstatic/",
        "^\/ready"
    ]
    custom_templates_dir = "/cnvrg-static"
    email_domains = [
    {{- range $_, $email := .EmailDomain }}
      "{{ $email }}",
    {{- end }}
    ]
    provider = "{{ .Provider }}"
    client_id = "{{ .ClientId }}"
    client_secret = "{{ .ClientSecret }}"
    redirect_url = "{{ .RedirectUrl }}"
    oidc_issuer_url = "{{ .OidcIssuerURL }}"
    scope = "{{ .Scope }}"
    cookie_name = "_cnvrg_auth"
    cookie_secure = false
    cookie_expire = "168h"
    cookie_httponly = true
    cookie_secret = "{{ randAlphaNum 32 }}"
    session_store_type = "redis"
    set_xauthrequest = true
    pass_access_token = true
    pass_authorization_header = true
    skip_jwt_bearer_tokens = true
    insecure_oidc_allow_unverified_email = {{ .InsecureOidcAllowUnverifiedEmail }}
    ssl_insecure_skip_verify  = {{ .SslInsecureSkipVerify }}
    whitelist_domains = [ "{{ .WhitelistDomain }}" ]
    cookie_domains = [ "{{ .CookieDomain }}" ]
    extra_jwt_issuers = [ "{{.ExtraJwtIssuer}}" ]
    redis_connection_idle_timeout = 5
    {{- $groupLen := len .Groups }}
    {{- if gt $groupLen 0 }}
    allowed_groups = [
    {{- range $_, $g := .Groups }}
      "{{ $g }}",
    {{- end }}
    ]
    {{- end }}
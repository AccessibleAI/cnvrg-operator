package v1

type SSO struct {
	Enabled            string `json:"enabled,omitempty"`
	Image              string `json:"image,omitempty"`
	RedisConnectionUrl string `json:"redisConnectionUrl,omitempty"`
	AdminUser          string `json:"adminUser,omitempty"`
	Provider           string `json:"provider,omitempty"`
	EmailDomain        string `json:"emailDomain,omitempty"`
	ClientID           string `json:"clientId,omitempty"`
	ClientSecret       string `json:"clientSecret,omitempty"`
	CookieSecret       string `json:"cookieSecret,omitempty"`
	AzureTenant        string `json:"azureTenant,omitempty"`
	OidcIssuerURL      string `json:"oidcIssuerUrl,omitempty"`
	KeysSecretRef      string `json:"keysSecretRef,omitempty"`
}

type OauthProxyServiceConf struct {
	SkipAuthRegex []string `json:"skipAuthRegex,omitempty"`
}

var ssoDefault = SSO{
	Enabled:            "false",
	Image:              "cnvrg/cnvrg-oauth-proxy:v7.0.1.c7",
	RedisConnectionUrl: "redis://redis:6379",
	KeysSecretRef:      "sso-keys",
	AdminUser:          "",
	Provider:           "",
	EmailDomain:        "",
	ClientID:           "",
	ClientSecret:       "",
	CookieSecret:       "",
	AzureTenant:        "", // if IDP is Azure AD
	OidcIssuerURL:      "", // if IDP oidc

}

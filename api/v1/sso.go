package v1

type SSO struct {
	Enabled       *bool    `json:"enabled,omitempty"`
	Image         string   `json:"image,omitempty"`
	AdminUser     string   `json:"adminUser,omitempty"`
	Provider      string   `json:"provider,omitempty"`
	EmailDomain   []string `json:"emailDomain,omitempty"`
	ClientID      string   `json:"clientId,omitempty"`
	ClientSecret  string   `json:"clientSecret,omitempty"`
	CookieSecret  string   `json:"cookieSecret,omitempty"`
	AzureTenant   string   `json:"azureTenant,omitempty"`
	OidcIssuerURL string   `json:"oidcIssuerUrl,omitempty"`
}

type OauthProxyServiceConf struct {
	SkipAuthRegex []string `json:"skipAuthRegex,omitempty"`
}

var ssoDefault = SSO{
	Enabled:       &defaultEnabled,
	Image:         "cnvrg/cnvrg-oauth-proxy:v7.0.1.c7",
	AdminUser:     "",
	Provider:      "",
	EmailDomain:   nil,
	ClientID:      "",
	ClientSecret:  "",
	CookieSecret:  "",
	AzureTenant:   "", // if IDP is Azure AD
	OidcIssuerURL: "", // if IDP oidc

}

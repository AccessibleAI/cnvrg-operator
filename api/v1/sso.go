package v1

type SSO struct {
	Enabled                          bool     `json:"enabled,omitempty"`
	Image                            string   `json:"image,omitempty"`
	AdminUser                        string   `json:"adminUser,omitempty"`
	Provider                         string   `json:"provider,omitempty"`
	EmailDomain                      []string `json:"emailDomain,omitempty"`
	ClientID                         string   `json:"clientId,omitempty"`
	ClientSecret                     string   `json:"clientSecret,omitempty"`
	CookieSecret                     string   `json:"cookieSecret,omitempty"`
	AzureTenant                      string   `json:"azureTenant,omitempty"`
	OidcIssuerURL                    string   `json:"oidcIssuerUrl,omitempty"`
	RealmName                        string   `json:"realmName,omitempty"`
	ServiceUrl                       string   `json:"serviceUrl,omitempty"`
	InsecureOidcAllowUnverifiedEmail bool     `json:"insecureOidcAllowUnverifiedEmail,omitempty"`
}

type OauthProxyServiceConf struct {
	SkipAuthRegex        []string `json:"skipAuthRegex,omitempty"`
	TokenValidationRegex []string `json:"tokenValidationRegex,omitempty"`
}

var ssoDefault = SSO{
	Enabled:                          false,
	Image:                            "saas-oauth2-proxy:latest",
	AdminUser:                        "",
	Provider:                         "",
	EmailDomain:                      nil,
	ClientID:                         "",
	ClientSecret:                     "",
	CookieSecret:                     "",
	AzureTenant:                      "", // if IDP is Azure AD
	OidcIssuerURL:                    "", // if IDP oidc
	RealmName:                        "",
	ServiceUrl:                       "",
	InsecureOidcAllowUnverifiedEmail: false,
}

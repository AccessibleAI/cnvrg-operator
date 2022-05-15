package v1

type SaaSSSO struct {
	Enabled         bool     `json:"enabled,omitempty"`
	ExtraJWTIssuers []string `json:"extraJWTIssuers,omitempty"`
	AllowedGroups   []string `json:"allowedGroups,omitempty"`
}

type SSO struct {
	Enabled                          bool     `json:"enabled,omitempty"`
	Image                            string   `json:"image,omitempty"`
	AdminUser                        string   `json:"adminUser,omitempty"`
	Provider                         string   `json:"provider,omitempty"`
	Scope                            string   `json:"scope"`
	EmailDomain                      []string `json:"emailDomain,omitempty"`
	ClientID                         string   `json:"clientId,omitempty"`
	ClientSecret                     string   `json:"clientSecret,omitempty"`
	CookieSecret                     string   `json:"cookieSecret,omitempty"`
	AzureTenant                      string   `json:"azureTenant,omitempty"`
	OidcIssuerURL                    string   `json:"oidcIssuerUrl,omitempty"`
	RealmName                        string   `json:"realmName,omitempty"`
	ServiceUrl                       string   `json:"serviceUrl,omitempty"`
	InsecureOidcAllowUnverifiedEmail bool     `json:"insecureOidcAllowUnverifiedEmail,omitempty"`
	SaaSSSO                          SaaSSSO  `json:"saaSSSO,omitempty"`
}

type OauthProxyServiceConf struct {
	SkipAuthRegex []string `json:"skipAuthRegex,omitempty"`
}

var ssoDefault = SSO{
	Enabled:                          false,
	Image:                            "oauth2-proxy:v7.2.0",
	AdminUser:                        "",
	Provider:                         "",
	Scope:                            "openid",
	EmailDomain:                      nil,
	ClientID:                         "",
	ClientSecret:                     "",
	CookieSecret:                     "",
	AzureTenant:                      "", // if IDP is Azure AD
	OidcIssuerURL:                    "", // if IDP oidc
	RealmName:                        "",
	ServiceUrl:                       "",
	InsecureOidcAllowUnverifiedEmail: false,
	SaaSSSO:                          SaaSSSO{Enabled: false},
}

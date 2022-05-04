package v1

type Pki struct {
	Enabled          bool   `json:"enabled,omitempty"`
	RootCaSecret     string `json:"rootCaSecret,omitempty"`
	PrivateKeySecret string `json:"privateKeySecret,omitempty"`
	PublicKeySecret  string `json:"publicKeySecret,omitempty"`
}

var pkiDefault = Pki{
	Enabled:          false,
	RootCaSecret:     "",
	PrivateKeySecret: "",
	PublicKeySecret:  "",
}

package v1

type Jwks struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

var jwksDefault = Jwks{
	Enabled: false,
	Image:   "cnvrg/cnvrg-jwks:v1",
}

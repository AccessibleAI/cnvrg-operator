package v1

type JwksCache struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Jwks struct {
	Enabled bool      `json:"enabled,omitempty"`
	Image   string    `json:"image,omitempty"`
	Cache   JwksCache `json:"cache"`
}

var jwksDefault = Jwks{
	Enabled: false,
	Image:   "cnvrg/jwks:latest",
	Cache: JwksCache{
		Enabled: true,
		Image:   "docker.io/redis",
	},
}

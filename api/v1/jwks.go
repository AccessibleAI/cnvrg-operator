package v1

type JwksCache struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Jwks struct {
	Enabled bool      `json:"enabled,omitempty"`
	Name    string    `json:"name,omitempty"`
	Image   string    `json:"image,omitempty"`
	Cache   JwksCache `json:"cache,omitempty"`
}

var jwksDefault = Jwks{
	Enabled: false,
	Image:   "cnvrg/jwks:latest",
	Name:    "cnvrg-jwks",
	Cache: JwksCache{
		Enabled: true,
		Image:   "docker.io/redis",
	},
}

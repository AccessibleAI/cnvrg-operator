package v1

type ConfigReloader struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

var defaultConfigReloader = ConfigReloader{
	Enabled: false,
	Image:   "config-reloader:latest",
}

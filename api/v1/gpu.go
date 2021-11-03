package v1

type NvidiaDp struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Gpu struct {
	NvidiaDp NvidiaDp `json:"nvidiaDp,omitempty"`
	HabanaDp HabanaDp `json:"habanaDp,omitempty"`
}

var nvidiaDpDefault = NvidiaDp{
	Enabled: false,
	Image:   "k8s-device-plugin:v0.9.0",
}

var gpuDefaults = Gpu{
	NvidiaDp: nvidiaDpDefault,
	HabanaDp: habanaDpDefault,
}

type HabanaDp struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

var habanaDpDefault = HabanaDp{
	Enabled: false,
	Image:   "vault.habana.ai/docker-k8s-device-plugin/docker-k8s-device-plugin:latest",
}

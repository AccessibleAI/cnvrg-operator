package v1

type NvidiaDp struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type HabanaDp struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type MetaGpuDp struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Gpu struct {
	NvidiaDp  NvidiaDp  `json:"nvidiaDp,omitempty"`
	HabanaDp  HabanaDp  `json:"habanaDp,omitempty"`
	MetaGpuDp MetaGpuDp `json:"metaGpuDp,omitempty"`
}

var nvidiaDpDefault = NvidiaDp{
	Enabled: false,
	Image:   "k8s-device-plugin:v0.9.0",
}

var habanaDpDefault = HabanaDp{
	Enabled: true,
	Image:   "vault.habana.ai/docker-k8s-device-plugin/docker-k8s-device-plugin:latest",
}

var metagpuDpDefaults = MetaGpuDp{
	Enabled: false,
	Image:   "metagpu-device-plugin:main",
}

var gpuDefaults = Gpu{
	NvidiaDp:  nvidiaDpDefault,
	HabanaDp:  habanaDpDefault,
	MetaGpuDp: metagpuDpDefaults,
}

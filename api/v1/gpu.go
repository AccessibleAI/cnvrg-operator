package v1

type NvidiaDp struct {
	Enabled string `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Gpu struct {
	NvidiaDp NvidiaDp `json:"nvidiaDp,omitempty"`
}

var nvidiaDpDefault = NvidiaDp{
	Enabled: "true",
	Image:   "nvcr.io/nvidia/k8s-device-plugin:v0.9.0",
}

var gpuDefaults = Gpu{
	NvidiaDp: nvidiaDpDefault,
}

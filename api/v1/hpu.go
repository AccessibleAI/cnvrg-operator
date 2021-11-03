package v1

type HpuDp struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Hpu struct {
	HpuDp HpuDp `json:"hpuDp,omitempty"`
}

var hpuDpDefault = HpuDp{
	Enabled: false,
	Image:   "vault.habana.ai/docker-k8s-device-plugin/docker-k8s-device-plugin:latest",
}

var hpuDefaults = Hpu{
	HpuDp: hpuDpDefault,
}

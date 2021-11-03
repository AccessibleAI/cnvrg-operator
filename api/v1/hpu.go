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
	Image:   "k8s-device-plugin:v0.9.0",
}

var hpuDefaults = Hpu{
	HpuDp: hpuDpDefault,
}

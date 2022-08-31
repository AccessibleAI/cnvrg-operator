package v1

type Istio struct {
	Enabled               bool              `json:"enabled,omitempty"`
	OperatorImage         string            `json:"operatorImage,omitempty"`
	PilotImage            string            `json:"pilotImage,omitempty"`
	ProxyImage            string            `json:"proxyImage,omitempty"`
	IngressSvcExtraPorts  []int             `json:"ingressSvcExtraPorts,omitempty"`
	ExternalIP            []string          `json:"externalIp,omitempty"`
	LBSourceRanges        []string          `json:"lbSourceRanges,omitempty"`
	IngressSvcAnnotations map[string]string `json:"ingressSvcAnnotations,omitempty"`
}

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

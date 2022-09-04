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

type NodeSelector struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type DcgmExporter struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type NvidiaDevicePlugin struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Nvidia struct {
	NodeSelector    NodeSelector       `json:"nodeSelector,omitempty"`
	DevicePlugin    NvidiaDevicePlugin `json:"devicePlugin,omitempty"`
	MetricsExporter DcgmExporter       `json:"metricsExporter,omitempty"`
}

type HabanaDevicePlugin struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type HabanaMetricsExporter struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Habana struct {
	DevicePlugin    HabanaDevicePlugin    `json:"devicePlugin,omitempty"`
	MetricsExporter HabanaMetricsExporter `json:"metricsExporter,omitempty"`
}

type Metagpu struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

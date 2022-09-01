package v1

var nvidiaDefaults = Nvidia{
	DevicePlugin: NvidiaDevicePlugin{
		Enabled: false,
		Image:   "k8s-device-plugin:v0.9.0",
	},
	MetricsExporter: DcgmExporter{},
}

var habanaDefaults = Habana{
	DevicePlugin: HabanaDevicePlugin{
		Enabled: true,
		Image:   "vault.habana.ai/docker-k8s-device-plugin/docker-k8s-device-plugin:latest",
	},
	MetricsExporter: HabanaMetricsExporter{},
}

var metagpuDefaults = Metagpu{
	Enabled: false,
	Image:   "metagpu-device-plugin:main",
}

var istioDefault = Istio{
	Enabled:               false,
	OperatorImage:         "istio-operator:1.10.2",
	PilotImage:            "pilot:1.10.2",
	ProxyImage:            "proxyv2:1.10.2",
	ExternalIP:            nil,
	IngressSvcAnnotations: nil,
	IngressSvcExtraPorts:  nil,
	LBSourceRanges:        nil,
}

var thirdPartyRegistryDefault = Registry{
	Name:     "cnvrg-third-party-registry",
	URL:      "docker.io",
	User:     "",
	Password: "",
}

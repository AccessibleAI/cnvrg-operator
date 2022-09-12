package v1

func DefaultCnvrgThirdPartySpec() CnvrgThirdPartySpec {
	return CnvrgThirdPartySpec{
		ImageHub: "docker.io/cnvrg",
		Istio: Istio{
			Enabled:               false,
			OperatorImage:         "istio-operator:1.10.2",
			PilotImage:            "pilot:1.10.2",
			ProxyImage:            "proxyv2:1.10.2",
			ExternalIP:            nil,
			IngressSvcAnnotations: nil,
			IngressSvcExtraPorts:  nil,
			LBSourceRanges:        nil,
		},
		Registry: Registry{
			Name:     "cnvrg-third-party-registry",
			URL:      "docker.io",
			User:     "",
			Password: "",
		},
		Nvidia: Nvidia{
			NodeSelector: NodeSelector{
				Key:   "accelerator",
				Value: "nvidia",
			},
			DevicePlugin: NvidiaDevicePlugin{
				Enabled: false,
				Image:   "k8s-device-plugin:v0.9.0",
			},
			MetricsExporter: DcgmExporter{
				Enabled: false,
				Image:   "nvcr.io/nvidia/k8s/dcgm-exporter:2.0.13-2.1.2-ubuntu20.04",
			},
		},
		Habana: Habana{
			DevicePlugin: HabanaDevicePlugin{
				Enabled: true,
				Image:   "vault.habana.ai/docker-k8s-device-plugin/docker-k8s-device-plugin:latest",
			},
			MetricsExporter: HabanaMetricsExporter{},
		},
		Metagpu: Metagpu{
			Enabled:      false,
			Image:        "metagpu-device-plugin:main",
			NodeSelector: map[string]string{},
		},
	}
}

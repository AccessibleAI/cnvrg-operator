package v1

type Istio struct {
	Enabled                  string `json:"enabled,omitempty"`
	OperatorImage            string `json:"operatorImage,omitempty"`
	Hub                      string `json:"hub,omitempty"`
	Tag                      string `json:"tag,omitempty"`
	ProxyImage               string `json:"proxyImage,omitempty"`
	MixerImage               string `json:"mixerImage,omitempty"`
	PilotImage               string `json:"pilotImage,omitempty"`
	GwName                   string `json:"gwName,omitempty"`
	ExternalIP               string `json:"externalIp,omitempty"`
	IngressSvcAnnotations    string `json:"ingressSvcAnnotations,omitempty"`
	IngressSvcExtraPorts     string `json:"ingressSvcExtraPorts,omitempty"`
	LoadBalancerSourceRanges string `json:"loadBalancerSourceRanges,omitempty"`
}

var istioDefault = Istio{
	Enabled:                  "false",
	OperatorImage:            "docker.io/istio/operator:1.8.1",
	Hub:                      "docker.io/istio",
	Tag:                      "1.8.1",
	ProxyImage:               "proxyv2",
	MixerImage:               "mixer",
	PilotImage:               "pilot",
	GwName:                   "cnvrg-gateway",
	ExternalIP:               "",
	IngressSvcAnnotations:    "",
	IngressSvcExtraPorts:     "",
	LoadBalancerSourceRanges: "",
}

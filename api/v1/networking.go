package v1

type Networking struct {
	Enabled     string  `json:"enabled"`
	IngressType string  `json:"ingressType"`
	HTTPS       HTTPS   `json:"https"`
	Istio       Istio   `json:"istio"`
	Ingress     Ingress `json:"ingress"`
}
type HTTPS struct {
	Enabled    string `json:"enabled"`
	Cert       string `json:"cert"`
	Key        string `json:"key"`
	CertSecret string `json:"certSecret"`
}
type Istio struct {
	Enabled                  string `json:"enabled"`
	OperatorImage            string `json:"operatorImage"`
	Hub                      string `json:"hub"`
	Tag                      string `json:"tag"`
	ProxyImage               string `json:"proxyImage"`
	MixerImage               string `json:"mixerImage"`
	PilotImage               string `json:"pilotImage"`
	GwName                   string `json:"gwName"`
	ExternalIP               string `json:"externalIp"`
	IngressSvcAnnotations    string `json:"ingressSvcAnnotations"`
	IngressSvcExtraPorts     string `json:"ingressSvcExtraPorts"`
	LoadBalancerSourceRanges string `json:"loadBalancerSourceRanges"`
}
type Ingress struct {
	Enabled         string `json:"enabled"`
	Timeout         string `json:"timeout"`
	RetriesAttempts int    `json:"retriesAttempts"`
	PerTryTimeout   string `json:"perTryTimeout"`
}

var networkingDefault = Networking{
	Enabled:     "true",
	IngressType: "istio", // openshift | istio | k8singress | nodeport
	HTTPS: HTTPS{
		Enabled:    "false",
		Cert:       "",
		Key:        "",
		CertSecret: "",
	},
	Istio: Istio{
		Enabled:                  "true",
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
	},
	Ingress: Ingress{
		Enabled:         "true",
		Timeout:         "18000s",
		RetriesAttempts: 5,
		PerTryTimeout:   "3600s",
	},
}

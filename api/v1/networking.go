package v1

// +kubebuilder:validation:Enum=istio;k8singress;openshift;nodeport
type IngressType string

const (
	IstioIngress     IngressType = "istio"
	NginxIngress     IngressType = "k8singress"
	OpenShiftIngress IngressType = "openshift"
	NodePortIngress  IngressType = "nodeport"
)

type Networking struct {
	Enabled     string      `json:"enabled,omitempty"`
	IngressType IngressType `json:"ingressType,omitempty"`
	HTTPS       HTTPS       `json:"https,omitempty"`
	Istio       Istio       `json:"istio,omitempty"`
	Ingress     Ingress     `json:"ingress,omitempty"`
}
type HTTPS struct {
	Enabled    string `json:"enabled,omitempty"`
	Cert       string `json:"cert,omitempty"`
	Key        string `json:"key,omitempty"`
	CertSecret string `json:"certSecret,omitempty"`
}
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
type Ingress struct {
	Enabled         string `json:"enabled,omitempty"`
	Timeout         string `json:"timeout,omitempty"`
	RetriesAttempts int    `json:"retriesAttempts,omitempty"`
	PerTryTimeout   string `json:"perTryTimeout,omitempty"`
}

var networkingDefault = Networking{
	Enabled:     "true",
	IngressType: IstioIngress,
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

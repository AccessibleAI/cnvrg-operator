package v1

// +kubebuilder:validation:Enum=istio;k8singress;openshift;nodeport
type IngressType string

const (
	IstioIngress     IngressType = "istio"
	NginxIngress     IngressType = "k8singress"
	OpenShiftIngress IngressType = "openshift"
	NodePortIngress  IngressType = "nodeport"
)

type Istio struct {
	Enabled                  *bool  `json:"enabled,omitempty"`
	OperatorImage            string `json:"operatorImage,omitempty"`
	Hub                      string `json:"hub,omitempty"`
	Tag                      string `json:"tag,omitempty"`
	ProxyImage               string `json:"proxyImage,omitempty"`
	MixerImage               string `json:"mixerImage,omitempty"`
	PilotImage               string `json:"pilotImage,omitempty"`
	ExternalIP               string `json:"externalIp,omitempty"`
	IngressSvcAnnotations    string `json:"ingressSvcAnnotations,omitempty"`
	IngressSvcExtraPorts     string `json:"ingressSvcExtraPorts,omitempty"`
	LoadBalancerSourceRanges string `json:"loadBalancerSourceRanges,omitempty"`
}

type Ingress struct {
	IngressType     IngressType `json:"ingressType,omitempty"`
	Timeout         string      `json:"timeout,omitempty"`
	RetriesAttempts int         `json:"retriesAttempts,omitempty"`
	PerTryTimeout   string      `json:"perTryTimeout,omitempty"`
	IstioGwName     string      `json:"istioGwName,omitempty"`
}

type HTTPS struct {
	Enabled    *bool  `json:"enabled,omitempty"`
	Cert       string `json:"cert,omitempty"`
	Key        string `json:"key,omitempty"`
	CertSecret string `json:"certSecret,omitempty"`
}

type CnvrgAppNetworking struct {
	Ingress Ingress `json:"ingress,omitempty"`
	HTTPS   HTTPS   `json:"https,omitempty"`
}

type CnvrgInfraNetworking struct {
	Ingress Ingress `json:"ingress,omitempty"`
	HTTPS   HTTPS   `json:"https,omitempty"`
	Istio   Istio   `json:"istio,omitempty"`
}

var istioDefault = Istio{
	Enabled:                  &defaultEnabled,
	OperatorImage:            "docker.io/istio/operator:1.8.5",
	Hub:                      "docker.io/istio",
	Tag:                      "1.8.5",
	ProxyImage:               "proxyv2",
	MixerImage:               "mixer",
	PilotImage:               "pilot",
	ExternalIP:               "",
	IngressSvcAnnotations:    "",
	IngressSvcExtraPorts:     "",
	LoadBalancerSourceRanges: "",
}

var httpsDefault = HTTPS{
	Enabled:    &defaultEnabled,
	Cert:       "",
	Key:        "",
	CertSecret: "",
}

var ingressAppDefault = Ingress{
	IngressType:     IstioIngress,
	Timeout:         "18000s",
	RetriesAttempts: 5,
	PerTryTimeout:   "3600s",
	IstioGwName:     "",
}

var ingressInfraDefault = Ingress{
	IngressType:     IstioIngress,
	Timeout:         "18000s",
	RetriesAttempts: 5,
	PerTryTimeout:   "3600s",
	IstioGwName:     "",
}

var cnvrgAppNetworkingDefault = CnvrgAppNetworking{
	Ingress: ingressAppDefault,
	HTTPS:   httpsDefault,
}

var cnvrgInfraNetworkingDefault = CnvrgInfraNetworking{
	Ingress: ingressInfraDefault,
	Istio:   istioDefault,
	HTTPS:   httpsDefault,
}

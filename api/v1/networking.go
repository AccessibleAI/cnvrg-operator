package v1

// +kubebuilder:validation:Enum=istio;ingress;openshift;nodeport
type IngressType string

var DoNotDeployIstioGwFlag string = "do-not-deploy"

const (
	IstioIngress     IngressType = "istio"
	NginxIngress     IngressType = "ingress"
	OpenShiftIngress IngressType = "openshift"
	NodePortIngress  IngressType = "nodeport"
)

type Istio struct {
	Enabled               *bool             `json:"enabled,omitempty"`
	Hub                   string            `json:"hub,omitempty"`
	Tag                   string            `json:"tag,omitempty"`
	IngressSvcExtraPorts  []int             `json:"ingressSvcExtraPorts,omitempty"`
	ExternalIP            []string          `json:"externalIp,omitempty"`
	LBSourceRanges        []string          `json:"lbSourceRanges,omitempty"`
	IngressSvcAnnotations map[string]string `json:"ingressSvcAnnotations,omitempty"`
}

type Ingress struct {
	Type            IngressType `json:"type,omitempty"`
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
	Enabled:               &defaultFalse,
	Hub:                   "docker.io/istio",
	Tag:                   "1.10.0",
	ExternalIP:            nil,
	IngressSvcAnnotations: nil,
	IngressSvcExtraPorts:  nil,
	LBSourceRanges:        nil,
}

var httpsDefault = HTTPS{
	Enabled:    &defaultFalse,
	Cert:       "",
	Key:        "",
	CertSecret: "",
}

var ingressAppDefault = Ingress{
	Type:            IstioIngress,
	Timeout:         "18000s",
	RetriesAttempts: 5,
	PerTryTimeout:   "3600s",
	IstioGwName:     "",
}

var ingressInfraDefault = Ingress{
	Type:            IstioIngress,
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

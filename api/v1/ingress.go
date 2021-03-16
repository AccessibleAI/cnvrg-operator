package v1

// +kubebuilder:validation:Enum=istio;k8singress;openshift;nodeport
type IngressType string

const (
	IstioIngress     IngressType = "istio"
	NginxIngress     IngressType = "k8singress"
	OpenShiftIngress IngressType = "openshift"
	NodePortIngress  IngressType = "nodeport"
)

type Ingress struct {
	Enabled         string      `json:"enabled,omitempty"`
	IngressType     IngressType `json:"ingressType,omitempty"`
	HTTPS           HTTPS       `json:"https,omitempty"`
	Timeout         string      `json:"timeout,omitempty"`
	RetriesAttempts int         `json:"retriesAttempts,omitempty"`
	PerTryTimeout   string      `json:"perTryTimeout,omitempty"`
	IstioGwName     string      `json:"istioGwName,omitempty"`
}
type HTTPS struct {
	Enabled    string `json:"enabled,omitempty"`
	Cert       string `json:"cert,omitempty"`
	Key        string `json:"key,omitempty"`
	CertSecret string `json:"certSecret,omitempty"`
}

var ingressDefault = Ingress{
	Enabled:         "true",
	IngressType:     IstioIngress,
	Timeout:         "18000s",
	RetriesAttempts: 5,
	PerTryTimeout:   "3600s",
	IstioGwName:     "cnvrg-gateway",
	HTTPS: HTTPS{
		Enabled:    "false",
		Cert:       "",
		Key:        "",
		CertSecret: "",
	},
}

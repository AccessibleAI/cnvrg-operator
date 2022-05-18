package v1

// +kubebuilder:validation:Enum=istio;ingress;openshift;nodeport
type IngressType string

var IstioGwName string = "istio-gw-%s"

const (
	IstioIngress     IngressType = "istio"
	NginxIngress     IngressType = "ingress"
	OpenShiftIngress IngressType = "openshift"
	NodePortIngress  IngressType = "nodeport"
)

type Istio struct {
	Enabled               bool              `json:"enabled,omitempty"`
	OperatorImage         string            `json:"operatorImage,omitempty"`
	PilotImage            string            `json:"pilotImage,omitempty"`
	ProxyImage            string            `json:"proxyImage,omitempty"`
	IngressSvcExtraPorts  []int             `json:"ingressSvcExtraPorts,omitempty"`
	ExternalIP            []string          `json:"externalIp,omitempty"`
	LBSourceRanges        []string          `json:"lbSourceRanges,omitempty"`
	IngressSvcAnnotations map[string]string `json:"ingressSvcAnnotations,omitempty"`
	EastWest              EastWest          `json:"eastWest,omitempty"`
}

type EastWest struct {
	Enabled        bool             `json:"enabled,omitempty"`
	Primary        bool             `json:"primary,omitempty"`
	ClusterName    string           `json:"clusterName,omitempty"`
	Network        string           `json:"network,omitempty"`
	MeshId         string           `json:"meshId,omitempty"`
	RemoteClusters []RemoteClusters `json:"remoteClusters,omitempty"`
}

type RemoteClusters struct {
	ClusterName string `json:"clusterName,omitempty"`
	Network     string `json:"network,omitempty"`
	MeshId      string `json:"meshId,omitempty"`
	Primary     bool   `json:"primary,omitempty"`
}

type Ingress struct {
	Type            IngressType `json:"type,omitempty"`
	Timeout         string      `json:"timeout,omitempty"`
	RetriesAttempts int         `json:"retriesAttempts,omitempty"`
	PerTryTimeout   string      `json:"perTryTimeout,omitempty"`
	IstioGwEnabled  bool        `json:"istioGwEnabled,omitempty"`
	IstioGwName     string      `json:"istioGwName,omitempty"`
}

type HTTPS struct {
	Enabled    bool   `json:"enabled,omitempty"`
	Cert       string `json:"cert,omitempty"`
	Key        string `json:"key,omitempty"`
	CertSecret string `json:"certSecret,omitempty"`
}

type CnvrgAppNetworking struct {
	Ingress Ingress `json:"ingress,omitempty"`
	HTTPS   HTTPS   `json:"https,omitempty"`
	Proxy   Proxy   `json:"proxy,omitempty"`
}

type CnvrgInfraNetworking struct {
	Ingress Ingress `json:"ingress,omitempty"`
	HTTPS   HTTPS   `json:"https,omitempty"`
	Istio   Istio   `json:"istio,omitempty"`
	Proxy   Proxy   `json:"proxy,omitempty"`
}

type Proxy struct {
	Enabled    bool     `json:"enabled,omitempty"`
	ConfigRef  string   `json:"configRef,omitempty"`
	HttpProxy  []string `json:"httpProxy,omitempty"`
	HttpsProxy []string `json:"httpsProxy,omitempty"`
	NoProxy    []string `json:"noProxy,omitempty"`
}

var defaultInfraProxy = Proxy{
	Enabled:    false,
	ConfigRef:  "infra-proxy",
	HttpProxy:  nil,
	HttpsProxy: nil,
	NoProxy:    nil,
}

var defaultAppProxy = Proxy{
	Enabled:    false,
	ConfigRef:  "cp-proxy",
	HttpProxy:  nil,
	HttpsProxy: nil,
	NoProxy:    nil,
}

var istioDefault = Istio{
	Enabled:               false,
	OperatorImage:         "istio-operator:1.13.3",
	PilotImage:            "pilot:1.13.3",
	ExternalIP:            nil,
	IngressSvcAnnotations: nil,
	IngressSvcExtraPorts:  nil,
	LBSourceRanges:        nil,
	EastWest: EastWest{
		Enabled:        false,
		Primary:        false,
		ClusterName:    "",
		Network:        "network1",
		MeshId:         "mesh1",
		RemoteClusters: nil,
	},
}

var httpsDefault = HTTPS{
	Enabled:    false,
	Cert:       "",
	Key:        "",
	CertSecret: "",
}

var ingressAppDefault = Ingress{
	Type:            IstioIngress,
	Timeout:         "18000s",
	RetriesAttempts: 5,
	PerTryTimeout:   "3600s",
	IstioGwEnabled:  false,
	IstioGwName:     "",
}

var ingressInfraDefault = Ingress{
	Type:            IstioIngress,
	Timeout:         "18000s",
	RetriesAttempts: 5,
	PerTryTimeout:   "3600s",
	IstioGwEnabled:  false,
	IstioGwName:     "",
}

var cnvrgAppNetworkingDefault = CnvrgAppNetworking{
	Ingress: ingressAppDefault,
	HTTPS:   httpsDefault,
	Proxy:   defaultAppProxy,
}

var cnvrgInfraNetworkingDefault = CnvrgInfraNetworking{
	Ingress: ingressInfraDefault,
	Istio:   istioDefault,
	HTTPS:   httpsDefault,
	Proxy:   defaultInfraProxy,
}

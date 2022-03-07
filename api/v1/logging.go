package v1

type AppInstance struct {
	SpecName string
	SpecNs   string
	EsUser   string
	EsPass   string
}

type Fluentbit struct {
	Enabled      bool              `json:"enabled,omitempty"`
	Image        string            `json:"image,omitempty"`
	Requests     Requests          `json:"requests,omitempty"`
	Limits       Limits            `json:"limits,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	LogsMounts   map[string]string `json:"logsMounts,omitempty"`
}

type Elastalert struct {
	Enabled        bool              `json:"enabled,omitempty"`
	Image          string            `json:"image,omitempty"`
	AuthProxyImage string            `json:"authProxyImage,omitempty"`
	CredsRef       string            `json:"credsRef,omitempty"`
	Port           int               `json:"port,omitempty"`
	NodePort       int               `json:"nodePort,omitempty"`
	StorageSize    string            `json:"storageSize,omitempty"`
	SvcName        string            `json:"svcName,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	Requests       Requests          `json:"requests,omitempty"`
	Limits         Limits            `json:"limits,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	PvcName        string            `json:"pvcName,omitempty"`
}

type Kibana struct {
	Enabled        bool                  `json:"enabled,omitempty"`
	ServiceAccount string                `json:"serviceAccount,omitempty"`
	SvcName        string                `json:"svcName,omitempty"`
	Port           int                   `json:"port,omitempty"`
	Image          string                `json:"image,omitempty"`
	NodePort       int                   `json:"nodePort,omitempty"`
	Requests       Requests              `json:"requests,omitempty"`
	Limits         Limits                `json:"limits,omitempty"`
	OauthProxy     OauthProxyServiceConf `json:"oauthProxy,omitempty"`
	CredsRef       string                `json:"credsRef,omitempty"`
}

type CnvrgAppLogging struct {
	Elastalert Elastalert `json:"elastalert,omitempty"`
	Kibana     Kibana     `json:"kibana,omitempty"`
}

type CnvrgInfraLogging struct {
	Fluentbit Fluentbit `json:"fluentbit,omitempty"`
}

var fluentbitDefault = Fluentbit{
	Enabled:      false,
	Image:        "cnvrg-fluentbit:v1.7.3",
	NodeSelector: nil,
	LogsMounts: map[string]string{
		"varlog":                 "/var/log",
		"varlibdockercontainers": "/var/lib/docker/containers",
	},
	Requests: Requests{
		Cpu:    "50m",
		Memory: "200Mi",
	},
	Limits: Limits{
		Cpu:    "2000m",
		Memory: "2Gi",
	},
}

var cnvrgAppLoggingDefault = CnvrgAppLogging{
	Elastalert: Elastalert{
		Enabled:        false,
		Image:          "elastalert:3.0.0-beta.1",
		CredsRef:       "elastalert-creds",
		AuthProxyImage: "nginx:1.20",
		Port:           80,
		NodePort:       32030,
		StorageSize:    "30Gi",
		SvcName:        "elastalert",
		StorageClass:   "",
		Requests: Requests{
			Cpu:    "100m",
			Memory: "200Mi",
		},
		Limits: Limits{
			Cpu:    "400m",
			Memory: "800Mi",
		},
		NodeSelector: nil,
		PvcName:      "elastalert-storage",
	},
	Kibana: Kibana{
		Enabled:        false,
		ServiceAccount: "kibana",
		SvcName:        "kibana",
		Port:           8080,
		Image:          "kibana-oss:7.8.1",
		NodePort:       30601,
		Requests: Requests{
			Cpu:    "100m",
			Memory: "200Mi",
		},
		Limits: Limits{
			Cpu:    "1000m",
			Memory: "2Gi",
		},
		OauthProxy: OauthProxyServiceConf{
			SkipAuthRegex:        nil,
			TokenValidationRegex: nil,
		},
		CredsRef: "kibana-creds",
	},
}

var cnvrgInfraLoggingDefault = CnvrgInfraLogging{
	Fluentbit: fluentbitDefault,
}

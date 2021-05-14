package v1

type AppInstance struct {
	SpecName string
	SpecNs   string
	EsUser   string
	EsPass   string
}

type Fluentbit struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Elastalert struct {
	Enabled      *bool             `json:"enabled,omitempty"`
	Image        string            `json:"image,omitempty"`
	Port         int               `json:"port,omitempty"`
	NodePort     int               `json:"nodePort,omitempty"`
	StorageSize  string            `json:"storageSize,omitempty"`
	SvcName      string            `json:"svcName,omitempty"`
	StorageClass string            `json:"storageClass,omitempty"`
	Requests     Requests          `json:"requests,omitempty"`
	Limits       Limits            `json:"limits,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

type Kibana struct {
	Enabled        *bool                 `json:"enabled,omitempty"`
	ServiceAccount string                `json:"serviceAccount,omitempty"`
	SvcName        string                `json:"svcName,omitempty"`
	Port           int                   `json:"port,omitempty"`
	Image          string                `json:"image,omitempty"`
	NodePort       int                   `json:"nodePort,omitempty"`
	Requests       Requests              `json:"requests,omitempty"`
	Limits         Limits                `json:"limits,omitempty"`
	OauthProxy     OauthProxyServiceConf `json:"oauthProxy,omitempty"`
}

type CnvrgAppLogging struct {
	Elastalert Elastalert `json:"elastalert,omitempty"`
	Kibana     Kibana     `json:"kibana,omitempty"`
}

type CnvrgInfraLogging struct {
	Fluentbit Fluentbit `json:"fluentbit,omitempty"`
}

var fluentbitDefault = Fluentbit{
	Enabled: &defaultEnabled,
	Image:   "cnvrg/cnvrg-fluentbit:v1.7.2",
}

var cnvrgAppLoggingDefault = CnvrgAppLogging{
	Elastalert: Elastalert{
		Enabled:      &defaultEnabled,
		Image:        "bitsensor/elastalert:3.0.0-beta.1",
		Port:         80,
		NodePort:     32030,
		StorageSize:  "30Gi",
		SvcName:      "elastalert",
		StorageClass: "",
		Requests: Requests{
			Cpu:    "100m",
			Memory: "200Mi",
		},
		Limits: Limits{
			Cpu:    "400m",
			Memory: "800Mi",
		},
		NodeSelector: nil,
	},
	Kibana: Kibana{
		Enabled:        &defaultEnabled,
		ServiceAccount: "kibana",
		SvcName:        "kibana",
		Port:           8080,
		Image:          "docker.elastic.co/kibana/kibana-oss:7.8.1",
		NodePort:       30601,
		Requests: Requests{
			Cpu:    "100m",
			Memory: "200Mi",
		},
		Limits: Limits{
			Cpu:    "1000m",
			Memory: "2Gi",
		},
		OauthProxy: OauthProxyServiceConf{SkipAuthRegex: nil},
	},
}

var cnvrgInfraLoggingDefault = CnvrgInfraLogging{
	Fluentbit: fluentbitDefault,
}

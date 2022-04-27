package v1

// +kubebuilder:validation:Enum=minio;aws;azure;gcp
type ObjectStorageType string

const (
	MinioObjectStorageType ObjectStorageType = "minio"
	AwsObjectStorageType   ObjectStorageType = "aws"
	AzureObjectStorageType ObjectStorageType = "azure"
	GcpObjectStorageType   ObjectStorageType = "gcp"
)

type ConsistentHash struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type SharedStorage struct {
	Enabled        bool           `json:"enabled,omitempty"`
	ConsistentHash ConsistentHash `json:"consistentHash,omitempty"`
}

type Limits struct {
	Cpu    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type Hpa struct {
	Enabled     bool `json:"enabled,omitempty"`
	Utilization int  `json:"utilization,omitempty"`
	MaxReplicas int  `json:"maxReplicas,omitempty"`
}

type Requests struct {
	Cpu    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type WebApp struct {
	Replicas                int                   `json:"replicas,omitempty"`
	Enabled                 bool                  `json:"enabled,omitempty"`
	Port                    int                   `json:"port,omitempty"`
	Requests                Requests              `json:"requests,omitempty"`
	Limits                  Limits                `json:"limits,omitempty"`
	SvcName                 string                `json:"svcName,omitempty"`
	NodePort                int                   `json:"nodePort,omitempty"`
	PassengerMaxPoolSize    int                   `json:"passengerMaxPoolSize,omitempty"`
	InitialDelaySeconds     int                   `json:"initialDelaySeconds,omitempty"`
	ReadinessPeriodSeconds  int                   `json:"readinessPeriodSeconds,omitempty"`
	ReadinessTimeoutSeconds int                   `json:"readinessTimeoutSeconds,omitempty"`
	FailureThreshold        int                   `json:"failureThreshold,omitempty"`
	OauthProxy              OauthProxyServiceConf `json:"oauthProxy,omitempty"`
	Hpa                     Hpa                   `json:"hpa,omitempty"`
}

type Sidekiq struct {
	Enabled  bool     `json:"enabled,omitempty"`
	Split    bool     `json:"split,omitempty"`
	Requests Requests `json:"requests,omitempty"`
	Limits   Limits   `json:"limits,omitempty"`
	Replicas int      `json:"replicas,omitempty"`
	Hpa      Hpa      `json:"hpa,omitempty"`
}

type Searchkiq struct {
	Enabled  bool     `json:"enabled,omitempty"`
	Requests Requests `json:"requests,omitempty"`
	Limits   Limits   `json:"limits,omitempty"`
	Replicas int      `json:"replicas,omitempty"`
	Hpa      Hpa      `json:"hpa,omitempty"`
}

type Systemkiq struct {
	Enabled  bool     `json:"enabled,omitempty"`
	Requests Requests `json:"requests,omitempty"`
	Limits   Limits   `json:"limits,omitempty"`
	Replicas int      `json:"replicas,omitempty"`
	Hpa      Hpa      `json:"hpa,omitempty"`
}

type CnvrgRouter struct {
	Enabled  bool   `json:"enabled,omitempty"`
	Image    string `json:"image,omitempty"`
	SvcName  string `json:"svcName,omitempty"`
	NodePort int    `json:"nodePort,omitempty"`
}

type Hyper struct {
	Enabled                 bool     `json:"enabled,omitempty"`
	Image                   string   `json:"image,omitempty"`
	Port                    int      `json:"port,omitempty"`
	Replicas                int      `json:"replicas,omitempty"`
	NodePort                int      `json:"nodePort,omitempty"`
	SvcName                 string   `json:"svcName,omitempty"`
	Token                   string   `json:"token,omitempty"`
	Requests                Requests `json:"requests,omitempty"`
	Limits                  Limits   `json:"limits,omitempty"`
	CPULimit                string   `json:"cpuLimit,omitempty"`
	MemoryLimit             string   `json:"memoryLimit,omitempty"`
	ReadinessPeriodSeconds  int      `json:"readinessPeriodSeconds,omitempty"`
	ReadinessTimeoutSeconds int      `json:"readinessTimeoutSeconds,omitempty"`
}

type CnvrgScheduler struct {
	Enabled  bool     `json:"enabled,omitempty"`
	Requests Requests `json:"requests,omitempty"`
	Limits   Limits   `json:"limits,omitempty"`
	Replicas int      `json:"replicas,omitempty"`
}

type CnvrgClusterProvisionerOperator struct {
	Enabled     bool     `json:"enabled,omitempty"`
	Requests    Requests `json:"requests,omitempty"`
	Limits      Limits   `json:"limits,omitempty"`
	Image       string   `json:"image,omitempty"`
	AwsCredsRef string   `json:"awsCredsRef,omitempty"`
}

type Registry struct {
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

type Ldap struct {
	Enabled       bool   `json:"enabled,omitempty"`
	Host          string `json:"host,omitempty"`
	Port          string `json:"port,omitempty"`
	Account       string `json:"account,omitempty"`
	Base          string `json:"base,omitempty"`
	AdminUser     string `json:"adminUser,omitempty"`
	AdminPassword string `json:"adminPassword,omitempty"`
	Ssl           string `json:"ssl,omitempty"`
}

type SMTP struct {
	Server            string `json:"server,omitempty"`
	Port              int    `json:"port,omitempty"`
	Username          string `json:"username,omitempty"`
	Password          string `json:"password,omitempty"`
	Domain            string `json:"domain,omitempty"`
	OpensslVerifyMode string `json:"opensslVerifyMode,omitempty"`
	Sender            string `json:"sender,omitempty"`
}

type ObjectStorage struct {
	Type             ObjectStorageType `json:"type,omitempty"`
	Bucket           string            `json:"bucket,omitempty"`
	Region           string            `json:"region,omitempty"`
	AccessKey        string            `json:"accessKey,omitempty"`
	SecretKey        string            `json:"secretKey,omitempty"`
	Endpoint         string            `json:"endpoint,omitempty"`
	AzureAccountName string            `json:"azureAccountName,omitempty"`
	AzureContainer   string            `json:"azureContainer,omitempty"`
	GcpProject       string            `json:"gcpProject,omitempty"`
	GcpSecretRef     string            `json:"gcpSecretRef,omitempty"`
}

type BaseConfig struct {
	JobsStorageClass     string            `json:"jobsStorageClass,omitempty"`
	FeatureFlags         map[string]string `json:"featureFlags,omitempty"`
	SentryURL            string            `json:"sentryUrl,omitempty"`
	RunJobsOnSelfCluster string            `json:"runJobsOnSelfCluster,omitempty"`
	AgentCustomTag       string            `json:"agentCustomTag,omitempty"`
	Intercom             string            `json:"intercom,omitempty"`
	CnvrgJobUID          string            `json:"cnvrgJobUid,omitempty"`
	CnvrgJobRbacStrict   bool              `json:"cnvrgJobRbacStrict,omitempty"`
	CnvrgPrivilegedJob   bool              `json:"cnvrgPrivilegedJob,omitempty"`
}

type ControlPlane struct {
	Image                           string                          `json:"image,omitempty"`
	WebApp                          WebApp                          `json:"webapp,omitempty"`
	Sidekiq                         Sidekiq                         `json:"sidekiq,omitempty"`
	Searchkiq                       Searchkiq                       `json:"searchkiq,omitempty"`
	Systemkiq                       Systemkiq                       `json:"systemkiq,omitempty"`
	Hyper                           Hyper                           `json:"hyper,omitempty"`
	CnvrgScheduler                  CnvrgScheduler                  `json:"cnvrgScheduler,omitempty"`
	CnvrgClusterProvisionerOperator CnvrgClusterProvisionerOperator `json:"cnvrgClusterProvisionerOperator,omitempty"`
	CnvrgRouter                     CnvrgRouter                     `json:"cnvrgRouter,omitempty"`
	BaseConfig                      BaseConfig                      `json:"baseConfig,omitempty"`
	Ldap                            Ldap                            `json:"ldap,omitempty"`
	SMTP                            SMTP                            `json:"smtp,omitempty"`
	ObjectStorage                   ObjectStorage                   `json:"objectStorage,omitempty"`
	Mpi                             Mpi                             `json:"mpi,omitempty"`
}

type Mpi struct {
	Enabled              bool              `json:"enabled,omitempty"`
	Image                string            `json:"image,omitempty"`
	KubectlDeliveryImage string            `json:"kubectlDeliveryImage,omitempty"`
	ExtraArgs            map[string]string `json:"extraArgs,omitempty"`
	Registry             Registry          `json:"registry,omitempty"`
	Requests             Requests          `json:"requests,omitempty"`
	Limits               Limits            `json:"limits,omitempty"`
}

var mpiDefault = Mpi{
	Enabled:              false,
	Image:                "mpioperator/mpi-operator:v0.2.3",
	KubectlDeliveryImage: "mpioperator/kubectl-delivery:v0.2.3",
	ExtraArgs:            nil,
	Requests: Requests{
		Cpu:    "100m",
		Memory: "100Mi",
	},
	Limits: Limits{
		Cpu:    "1000m",
		Memory: "1Gi",
	},
	Registry: Registry{
		Name:     "mpi-private-registry",
		URL:      "docker.io",
		User:     "",
		Password: "",
	},
}

var hpa = Hpa{
	Enabled:     false,
	Utilization: 85,
	MaxReplicas: 5,
}

var appRegistryDefault = Registry{

	Name:     "cnvrg-app-registry",
	URL:      "docker.io",
	User:     "",
	Password: "",
}

var infraRegistryDefault = Registry{

	Name:     "cnvrg-infra-registry",
	URL:      "docker.io",
	User:     "",
	Password: "",
}

var controlPlaneDefault = ControlPlane{
	Image: "core:3.6.99",

	WebApp: WebApp{
		Enabled:  false,
		Replicas: 1,
		Port:     8080,
		Requests: Requests{
			Cpu:    "500m",
			Memory: "4Gi",
		},
		Limits: Limits{
			Cpu:    "4",
			Memory: "8Gi",
		},
		SvcName:                 "app",
		NodePort:                30080,
		PassengerMaxPoolSize:    50,
		InitialDelaySeconds:     10,
		ReadinessPeriodSeconds:  25,
		ReadinessTimeoutSeconds: 20,
		FailureThreshold:        5,
		OauthProxy: OauthProxyServiceConf{
			SkipAuthRegex: []string{
				`\/assets`,
				`\/healthz`,
				`\/public`,
				`\/pack`,
				`\/vscode.tar.gz`,
				`\/jupyter.vsix`,
				`\/gitlens.vsix`,
				`\/ms-python-release.vsix`,
				`\/webhooks`,
				`\/api/v2/metrics`,
				`\/api/v1/events/endpoint_rule_alert`,
			},
			TokenValidationRegex: []string{
				`^\/api`,
			},
		},
		Hpa: hpa,
	},

	Sidekiq: Sidekiq{
		Enabled: false,
		Split:   false,
		Requests: Requests{
			Cpu:    "200m",
			Memory: "3750Mi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "8Gi",
		},
		Replicas: 2,
		Hpa:      hpa,
	},

	Searchkiq: Searchkiq{
		Enabled: false,
		Requests: Requests{
			Cpu:    "200m",
			Memory: "1Gi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "8Gi",
		},
		Replicas: 1,
		Hpa:      hpa,
	},

	Systemkiq: Systemkiq{
		Enabled: false,
		Requests: Requests{
			Cpu:    "300m",
			Memory: "2Gi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "8Gi",
		},
		Replicas: 1,
		Hpa:      hpa,
	},

	Hyper: Hyper{
		Enabled:  false,
		Image:    "hyper-server:latest",
		Port:     5050,
		Replicas: 1,
		NodePort: 30050,
		SvcName:  "hyper",
		Token:    "token",
		Requests: Requests{
			Cpu:    "100m",
			Memory: "200Mi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "4Gi",
		},
		ReadinessPeriodSeconds:  100,
		ReadinessTimeoutSeconds: 60,
	},

	CnvrgScheduler: CnvrgScheduler{
		Enabled: false,
		Requests: Requests{
			Cpu:    "200m",
			Memory: "1000Mi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "4Gi",
		},
		Replicas: 1,
	},

	CnvrgClusterProvisionerOperator: CnvrgClusterProvisionerOperator{
		Enabled: false,
		Requests: Requests{
			Cpu:    "200m",
			Memory: "1Gi",
		},
		Limits: Limits{
			Cpu:    "2",
			Memory: "4Gi",
		},
		Image:       "cnvrg/ccp-operator:v1",
		AwsCredsRef: "",
	},

	CnvrgRouter: CnvrgRouter{
		Enabled:  false,
		Image:    "nginx:1.21.0",
		SvcName:  "cnvrg-router",
		NodePort: 30081,
	},

	Mpi: mpiDefault,

	BaseConfig: BaseConfig{
		JobsStorageClass:   "",
		FeatureFlags:       nil,
		SentryURL:          "",
		AgentCustomTag:     "latest",
		Intercom:           "true",
		CnvrgJobUID:        "1000",
		CnvrgJobRbacStrict: false,
		CnvrgPrivilegedJob: true,
	},

	ObjectStorage: ObjectStorage{

		Type:             MinioObjectStorageType,
		Bucket:           "cnvrg-storage",
		AccessKey:        "",
		SecretKey:        "",
		Endpoint:         "",
		AzureAccountName: "",
		AzureContainer:   "",
		Region:           "eastus",
		GcpProject:       "",
		GcpSecretRef:     "gcp-storage-secret",
	},

	Ldap: Ldap{
		Enabled:       false,
		Host:          "",
		Port:          "",
		Account:       "userPrincipalName",
		Base:          "", // dc=my-domain,dc=local
		AdminUser:     "",
		AdminPassword: "",
		Ssl:           "", // true/false
	},

	SMTP: SMTP{
		Server:            "",
		Port:              587,
		Username:          "",
		Password:          "",
		Domain:            "",
		OpensslVerifyMode: "",
		Sender:            "info@cnvrg.io",
	},
}

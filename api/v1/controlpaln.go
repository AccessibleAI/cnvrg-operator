package v1

type WebApp struct {
	Replicas                int    `json:"replicas,omitempty"`
	Enabled                 string `json:"enabled,omitempty"`
	Image                   string `json:"image,omitempty"`
	Port                    int    `json:"port,omitempty"`
	CPU                     int    `json:"cpu,omitempty"`
	Memory                  string `json:"memory,omitempty"`
	SvcName                 string `json:"svcName,omitempty"`
	NodePort                int    `json:"nodePort,omitempty"`
	PassengerMaxPoolSize    int    `json:"passengerMaxPoolSize,omitempty"`
	InitialDelaySeconds     int    `json:"initialDelaySeconds,omitempty"`
	ReadinessPeriodSeconds  int    `json:"readinessPeriodSeconds,omitempty"`
	ReadinessTimeoutSeconds int    `json:"readinessTimeoutSeconds,omitempty"`
	FailureThreshold        int    `json:"failureThreshold,omitempty"`
}

type Sidekiq struct {
	Enabled     string `json:"enabled,omitempty"`
	Split       string `json:"split,omitempty"`
	CPU         string `json:"cpu,omitempty"`
	Memory      string `json:"memory,omitempty"`
	Replicas    int    `json:"replicas,omitempty"`
	KillTimeout int    `json:"killTimeout,omitempty"`
}
type Searchkiq struct {
	Enabled     string `json:"enabled,omitempty"`
	CPU         string `json:"cpu,omitempty"`
	Memory      string `json:"memory,omitempty"`
	Replicas    int    `json:"replicas,omitempty"`
	KillTimeout int    `json:"killTimeout,omitempty"`
}
type Systemkiq struct {
	Enabled     string `json:"enabled,omitempty"`
	CPU         string `json:"cpu,omitempty"`
	Memory      string `json:"memory,omitempty"`
	Replicas    int    `json:"replicas,omitempty"`
	KillTimeout int    `json:"killTimeout,omitempty"`
}
type Registry struct {
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}
type Hyper struct {
	Enabled                 string `json:"enabled,omitempty"`
	Image                   string `json:"image,omitempty"`
	Port                    int    `json:"port,omitempty"`
	Replicas                int    `json:"replicas,omitempty"`
	NodePort                int    `json:"nodePort,omitempty"`
	SvcName                 string `json:"svcName,omitempty"`
	Token                   string `json:"token,omitempty"`
	CPURequest              string `json:"cpuRequest,omitempty"`
	MemoryRequest           string `json:"memoryRequest,omitempty"`
	CPULimit                int    `json:"cpuLimit,omitempty"`
	MemoryLimit             string `json:"memoryLimit,omitempty"`
	EnableReadinessProbe    string `json:"enableReadinessProbe,omitempty"`
	ReadinessPeriodSeconds  int    `json:"readinessPeriodSeconds,omitempty"`
	ReadinessTimeoutSeconds int    `json:"readinessTimeoutSeconds,omitempty"`
}
type Seeder struct {
	Image           string `json:"image,omitempty"`
	SeedCmd         string `json:"seedCmd,omitempty"`
	CreateBucketCmd string `json:"createBucketCmd,omitempty"`
}
type Ldap struct {
	Enabled       string `json:"enabled,omitempty"`
	Host          string `json:"host,omitempty"`
	Port          string `json:"port,omitempty"`
	Account       string `json:"account,omitempty"`
	Base          string `json:"base,omitempty"`
	AdminUser     string `json:"adminUser,omitempty"`
	AdminPassword string `json:"adminPassword,omitempty"`
	Ssl           string `json:"ssl,omitempty"`
}
type Rbac struct {
	Role               string `json:"role,omitempty"`
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	RoleBindingName    string `json:"roleBindingName,omitempty"`
}
type SMTP struct {
	Server   string `json:"server,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Domain   string `json:"domain,omitempty"`
}

type OauthProxy struct {
	Enabled       string   `json:"enabled"`
	Image         string   `json:"image"`
	AdminUser     string   `json:"adminUser"`
	Provider      string   `json:"provider"`
	EmailDomain   string   `json:"emailDomain"`
	RedirectURI   string   `json:"redirectUri"`
	ClientID      string   `json:"clientId"`
	ClientSecret  string   `json:"clientSecret"`
	CookieSecret  string   `json:"cookieSecret"`
	AzureTenant   string   `json:"azureTenant"`
	OidcIssuerURL string   `json:"oidcIssuerUrl"`
	SkipAuthRegex []string `json:"skipAuthRegex"`
}

type ObjectStorage struct {
	CnvrgStorageType             string `json:"cnvrgStorageType,omitempty"`
	CnvrgStorageBucket           string `json:"cnvrgStorageBucket,omitempty"`
	CnvrgStorageAccessKey        string `json:"cnvrgStorageAccessKey,omitempty"`
	CnvrgStorageSecretKey        string `json:"cnvrgStorageSecretKey,omitempty"`
	CnvrgStorageEndpoint         string `json:"cnvrgStorageEndpoint,omitempty"`
	MinioSseMasterKey            string `json:"minioSseMasterKey,omitempty"`
	CnvrgStorageAzureAccessKey   string `json:"cnvrgStorageAzureAccessKey,omitempty"`
	CnvrgStorageAzureAccountName string `json:"cnvrgStorageAzureAccountName,omitempty"`
	CnvrgStorageAzureContainer   string `json:"cnvrgStorageAzureContainer,omitempty"`
	CnvrgStorageRegion           string `json:"cnvrgStorageRegion,omitempty"`
	CnvrgStorageProject          string `json:"cnvrgStorageProject,omitempty"`
	GcpStorageSecret             string `json:"gcpStorageSecret,omitempty"`
	GcpKeyfileMountPath          string `json:"gcpKeyfileMountPath,omitempty"`
	GcpKeyfileName               string `json:"gcpKeyfileName,omitempty"`
	SecretKeyBase                string `json:"secretKeyBase,omitempty"`
	StsIv                        string `json:"stsIv,omitempty"`
	StsKey                       string `json:"stsKey,omitempty"`
}

type BaseConfig struct {
	JobsStorageClass     string            `json:"jobsStorageClass,omitempty"`
	FeatureFlags         map[string]string `json:"featureFlags,omitempty"`
	SentryURL            string            `json:"sentryUrl,omitempty"`
	PassengerAppEnv      string            `json:"passengerAppEnv,omitempty"`
	RailsEnv             string            `json:"railsEnv,omitempty"`
	RunJobsOnSelfCluster string            `json:"runJobsOnSelfCluster,omitempty"`
	DefaultComputeConfig string            `json:"defaultComputeConfig,omitempty"`
	DefaultComputeName   string            `json:"defaultComputeName,omitempty"`
	UseStdout            string            `json:"useStdout,omitempty"`
	ExtractTagsFromCmd   string            `json:"extractTagsFromCmd,omitempty"`
	CheckJobExpiration   string            `json:"checkJobExpiration,omitempty"`
	AgentCustomTag       string            `json:"agentCustomTag,omitempty"`
	Intercom             string            `json:"intercom,omitempty"`
	CnvrgJobUID          string            `json:"cnvrgJobUid,omitempty"`
}
type CnvrgRouter struct {
	Enabled  string `json:"enabled,omitempty"`
	Image    string `json:"image,omitempty"`
	SvcName  string `json:"svcName,omitempty"`
	NodePort int    `json:"nodePort,omitempty"`
	Port     int    `json:"port,omitempty"`
}

type ControlPlan struct {
	WebApp        WebApp        `json:"webapp,omitempty"`
	Sidekiq       Sidekiq       `json:"sidekiq,omitempty"`
	Searchkiq     Searchkiq     `json:"searchkiq,omitempty"`
	Systemkiq     Systemkiq     `json:"systemkiq,omitempty"`
	CnvrgRouter   CnvrgRouter   `json:"cnvrgRouter,omitempty"`
	Hyper         Hyper         `json:"hyper,omitempty"`
	Seeder        Seeder        `json:"seeder,omitempty"`
	BaseConfig    BaseConfig    `json:"baseConfig,omitempty"`
	Ldap          Ldap          `json:"ldap,omitempty"`
	Registry      Registry      `json:"registry,omitempty"`
	Rbac          Rbac          `json:"rbac,omitempty"`
	SMTP          SMTP          `json:"smtp,omitempty"`
	Tenancy       Tenancy       `json:"tenancy,omitempty"`
	OauthProxy    OauthProxy    `json:"oauthProxy,omitempty"`
	ObjectStorage ObjectStorage `json:"objectStorage,omitempty"`
}

type Tenancy struct {
	Enabled        string `json:"enabled,omitempty"`
	DedicatedNodes string `json:"dedicatedNodes,omitempty"`
	Key            string `json:"key,omitempty"`
	Value          string `json:"value,omitempty"`
}

type Cnvrg struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var controlPlanDefault = ControlPlan{

	WebApp: WebApp{
		Replicas:                1,
		Enabled:                 "true",
		Image:                   "cnvrg/core:3.1.5",
		Port:                    80,
		CPU:                     2,
		Memory:                  "4Gi",
		SvcName:                 "app",
		NodePort:                30080,
		PassengerMaxPoolSize:    20,
		InitialDelaySeconds:     10,
		ReadinessPeriodSeconds:  25,
		ReadinessTimeoutSeconds: 20,
		FailureThreshold:        4,
	},

	Sidekiq: Sidekiq{
		Enabled:     "true",
		Split:       "true",
		CPU:         "1750m",
		Memory:      "3750Mi",
		Replicas:    2,
		KillTimeout: 60,
	},

	Searchkiq: Searchkiq{
		Enabled:     "true",
		CPU:         "750m",
		Memory:      "750Mi",
		Replicas:    1,
		KillTimeout: 60,
	},

	Systemkiq: Systemkiq{
		Enabled:     "false",
		CPU:         "500m",
		Memory:      "500Mi",
		Replicas:    1,
		KillTimeout: 60,
	},

	Hyper: Hyper{
		Enabled:                 "true",
		Image:                   "cnvrg/hyper-server:latest",
		Port:                    5050,
		Replicas:                1,
		NodePort:                30050,
		SvcName:                 "hyper",
		Token:                   "token",
		CPURequest:              "100m",
		MemoryRequest:           "200Mi",
		CPULimit:                2,
		MemoryLimit:             "4Gi",
		EnableReadinessProbe:    "true",
		ReadinessPeriodSeconds:  100,
		ReadinessTimeoutSeconds: 60,
	},

	Seeder: Seeder{
		Image:           "docker.io/cnvrg/cnvrg-boot:v0.26-tenancy",
		SeedCmd:         "rails db:migrate && rails db:seed && rails libraries:update",
		CreateBucketCmd: "mb.sh",
	},

	CnvrgRouter: CnvrgRouter{
		Enabled:  "false",
		Image:    "nginx",
		SvcName:  "routing-service",
		NodePort: 30081,
		Port:     80,
	},

	BaseConfig: BaseConfig{
		JobsStorageClass:     "",
		FeatureFlags:         map[string]string{},
		SentryURL:            "https://4409141e4a204282bd1f5c021e587509:dc15f684faa9479a839cf913b98b4ee2@sentry.cnvrg.io/32",
		PassengerAppEnv:      "app",
		RailsEnv:             "app",
		RunJobsOnSelfCluster: "true",
		DefaultComputeConfig: "/opt/kube",
		DefaultComputeName:   "default",
		UseStdout:            "true",
		ExtractTagsFromCmd:   "false",
		CheckJobExpiration:   "true",
		AgentCustomTag:       "latest",
		Intercom:             "true",
		CnvrgJobUID:          "1000",
	},

	ObjectStorage: ObjectStorage{
		CnvrgStorageType:             "minio",
		CnvrgStorageBucket:           "cnvrg-storage",
		CnvrgStorageAccessKey:        "AKIAIOSFODNN7EXAMPLE",
		CnvrgStorageSecretKey:        "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		CnvrgStorageEndpoint:         "",
		MinioSseMasterKey:            "my-minio-key:a310aadcefdb634b748ae31225f175e3f64591f955dfc66ccc20e128a6817ff9",
		CnvrgStorageAzureAccessKey:   "",
		CnvrgStorageAzureAccountName: "",
		CnvrgStorageAzureContainer:   "",
		CnvrgStorageRegion:           "eastus",
		CnvrgStorageProject:          "",
		SecretKeyBase:                "0d2b33c2cc19cfaa838d3c354354a18fcc92beaaa8e97889ef99341c8aaf963ad3afcf0f7c20454cabb5c573c3fc35b60221034e109f4fb651ed1415bf61e9d5",
		StsIv:                        "DeJ/CGz/Hkb/IbRe4t1xLg==",
		StsKey:                       "05646d3cbf8baa5be7150b4283eda07d",
		GcpStorageSecret:             "gcp-storage-secret",
		GcpKeyfileMountPath:          "/tmp/gcp_keyfile",
		GcpKeyfileName:               "key.json",
	},

	Ldap: Ldap{
		Enabled:       "false",
		Host:          "",
		Port:          "",
		Account:       "userPrincipalName",
		Base:          "", // dc=my-domain,dc=local
		AdminUser:     "",
		AdminPassword: "",
		Ssl:           "", // true/false
	},

	Registry: Registry{
		Name:     "cnvrg-registry",
		URL:      "docker.io",
		User:     "",
		Password: "",
	},

	Rbac: Rbac{
		Role:               "cnvrg-control-plan-role",
		ServiceAccountName: "cnvrg",
		RoleBindingName:    "cnvrg-control-plan-binding",
	},

	SMTP: SMTP{
		Server:   "",
		Port:     "",
		Username: "",
		Password: "",
		Domain:   "",
	},

	Tenancy: Tenancy{
		Enabled:        "false",
		DedicatedNodes: "false",
		Key:            "cnvrg-taint",
		Value:          "true",
	},

	OauthProxy: OauthProxy{
		Enabled:       "false",
		Image:         "cnvrg/cnvrg-oauth-proxy:v7.0.1.c1",
		AdminUser:     "",
		Provider:      "",
		EmailDomain:   "",
		RedirectURI:   "",
		ClientID:      "",
		ClientSecret:  "",
		CookieSecret:  "",
		AzureTenant:   "", // if IDP is Azure AD
		OidcIssuerURL: "", // if IDP oidc
		SkipAuthRegex: []string{
			`^\/api`,
			`^\/oauth/`,
			`\/assets`,
			`\/healthz`,
			`\/public`,
			`\/pack`,
			`\/vscode.tar.gz`,
			`\/gitlens.vsix`,
			`\/ms-python-release.vsix`,
		},
	},
}

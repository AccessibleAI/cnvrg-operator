package v1

type ConsistentHash struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
type SharedStorage struct {
	Enabled          string         `json:"enabled,omitempty"`
	UseExistingClaim string         `json:"useExistingClaim,omitempty"`
	ConsistentHash   ConsistentHash `json:"consistentHash,omitempty"`
}
type Minio struct {
	Enabled       string        `json:"enabled,omitempty"`
	Replicas      int           `json:"replicas,omitempty"`
	Image         string        `json:"image,omitempty"`
	Port          int           `json:"port,omitempty"`
	StorageSize   string        `json:"storageSize,omitempty"`
	SvcName       string        `json:"svcName,omitempty"`
	NodePort      int           `json:"nodePort,omitempty"`
	StorageClass  string        `json:"storageClass,omitempty"`
	CPURequest    int           `json:"cpuRequest,omitempty"`
	MemoryRequest string        `json:"memoryRequest,omitempty"`
	SharedStorage SharedStorage `json:"sharedStorage,omitempty"`
}

var minioDefaults = Minio{
	Enabled:       "true",
	Replicas:      1,
	Image:         "docker.io/minio/minio:RELEASE.2020-09-17T04-49-20Z",
	Port:          9000,
	StorageSize:   "100Gi",
	SvcName:       "minio",
	NodePort:      30090,
	StorageClass:  "use-default",
	CPURequest:    1,
	MemoryRequest: "2Gi",
	SharedStorage: SharedStorage{
		Enabled:          "enabled",
		UseExistingClaim: "",
		ConsistentHash: ConsistentHash{
			Key:   "httpQueryParameterName",
			Value: "uploadId",
		},
	},
}

type Limits struct {
	CPU    int    `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}
type Requests struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}
type Redis struct {
	Enabled      string   `json:"enabled,omitempty"`
	Image        string   `json:"image,omitempty"`
	SvcName      string   `json:"svcName,omitempty"`
	Port         int      `json:"port,omitempty"`
	Appendonly   string   `json:"appendonly,omitempty"`
	StorageSize  string   `json:"storageSize,omitempty"`
	StorageClass string   `json:"storageClass,omitempty"`
	Limits       Limits   `json:"limits,omitempty"`
	Requests     Requests `json:"requests,omitempty"`
}

var redisDefault = Redis{
	Enabled:      "true",
	Image:        "docker.io/cnvrg/cnvrg-redis:v3.0.5.c2",
	SvcName:      "redis",
	Port:         6379,
	Appendonly:   "yes",
	StorageSize:  "10Gi",
	StorageClass: "use-default",
	Limits: Limits{
		CPU:    1,
		Memory: "2Gi",
	},
	Requests: Requests{
		CPU:    "500m",
		Memory: "1Gi",
	},
}

type HugePages struct {
	Enabled string `json:"enabled,omitempty"`
	Size    string `json:"size,omitempty"`
	Memory  string `json:"memory,omitempty"`
}

type Pg struct {
	Enabled        string    `json:"enabled,omitempty"`
	SecretName     string    `json:"secretName,omitempty"`
	Image          string    `json:"image,omitempty"`
	Port           int       `json:"port,omitempty"`
	StorageSize    string    `json:"storageSize,omitempty"`
	SvcName        string    `json:"svcName,omitempty"`
	Dbname         string    `json:"dbname,omitempty"`
	Pass           string    `json:"pass,omitempty"`
	User           string    `json:"user,omitempty"`
	RunAsUser      int       `json:"runAsUser,omitempty"`
	FsGroup        int       `json:"fsGroup,omitempty"`
	StorageClass   string    `json:"storageClass,omitempty"`
	CPURequest     int       `json:"cpuRequest,omitempty"`
	MemoryRequest  string    `json:"memoryRequest,omitempty"`
	MaxConnections int       `json:"maxConnections,omitempty"`
	SharedBuffers  string    `json:"sharedBuffers,omitempty"`
	HugePages      HugePages `json:"hugePages,omitempty"`
	Fixpg          string    `json:"fixpg,omitempty"`
}

var pgDefault = Pg{
	Enabled:        "true",
	SecretName:     "cnvrg-pg-secret",
	Image:          "centos/postgresql-12-centos7",
	Port:           5432,
	StorageSize:    "80Gi",
	SvcName:        "postgres",
	Dbname:         "cnvrg_production",
	Pass:           "pg_pass",
	User:           "cnvrg",
	RunAsUser:      26,
	FsGroup:        26,
	StorageClass:   "use-default",
	CPURequest:     4,
	MemoryRequest:  "4Gi",
	MaxConnections: 100,
	SharedBuffers:  "64MB",
	Fixpg:          "true",
	HugePages: HugePages{
		Enabled: "false",
		Size:    "2Mi",
		Memory:  "",
	},
}

type WebApp struct {
	Replicas                int                   `json:"replicas,omitempty"`
	Enabled                 string                `json:"enabled,omitempty"`
	Image                   string                `json:"image,omitempty"`
	Port                    int                   `json:"port,omitempty"`
	CPU                     int                   `json:"cpu,omitempty"`
	Memory                  string                `json:"memory,omitempty"`
	SvcName                 string                `json:"svcName,omitempty"`
	NodePort                int                   `json:"nodePort,omitempty"`
	PassengerMaxPoolSize    int                   `json:"passengerMaxPoolSize,omitempty"`
	InitialDelaySeconds     int                   `json:"initialDelaySeconds,omitempty"`
	ReadinessPeriodSeconds  int                   `json:"readinessPeriodSeconds,omitempty"`
	ReadinessTimeoutSeconds int                   `json:"readinessTimeoutSeconds,omitempty"`
	FailureThreshold        int                   `json:"failureThreshold,omitempty"`
	OauthProxy              OauthProxyServiceConf `json:"oauthProxy"`
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
	CcpStorageClass      string            `json:"ccpStorageClass,omitempty"`
	HostpathNode         string            `json:"hostpathNode,omitempty"`
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
	ObjectStorage ObjectStorage `json:"objectStorage,omitempty"`
	Pg            Pg            `json:"pg,omitempty"`
	Minio         Minio         `json:"minio,omitempty"`
	Redis         Redis         `json:"redis,omitempty"`
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

var registryDefault = Registry{
	Name:     "cnvrg-registry",
	URL:      "docker.io",
	User:     "",
	Password: "",
}

var controlPlanDefault = ControlPlan{
	Pg:    pgDefault,
	Minio: minioDefaults,
	Redis: redisDefault,

	WebApp: WebApp{
		Replicas:                1,
		Enabled:                 "true",
		Image:                   "cnvrg/core:3.1.5",
		Port:                    8080,
		CPU:                     2,
		Memory:                  "4Gi",
		SvcName:                 "app",
		NodePort:                30080,
		PassengerMaxPoolSize:    20,
		InitialDelaySeconds:     10,
		ReadinessPeriodSeconds:  25,
		ReadinessTimeoutSeconds: 20,
		FailureThreshold:        4,
		OauthProxy: OauthProxyServiceConf{
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
				`^\/api\/health`,
			},
		},
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
		FeatureFlags:         nil,
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
		CcpStorageClass:      "",
		HostpathNode:         "",
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

	Registry: registryDefault,

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
}

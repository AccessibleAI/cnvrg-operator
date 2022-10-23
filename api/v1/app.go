package v1

// +kubebuilder:validation:Enum=minio;aws;azure;gcp
type ObjectStorageType string

const (
	MinioObjectStorageType ObjectStorageType = "minio"
	AwsObjectStorageType   ObjectStorageType = "aws"
	AzureObjectStorageType ObjectStorageType = "azure"
	GcpObjectStorageType   ObjectStorageType = "gcp"
)

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
	Replicas                int      `json:"replicas,omitempty"`
	Enabled                 bool     `json:"enabled,omitempty"`
	Port                    int      `json:"port,omitempty"`
	Requests                Requests `json:"requests,omitempty"`
	Limits                  Limits   `json:"limits,omitempty"`
	SvcName                 string   `json:"svcName,omitempty"`
	NodePort                int      `json:"nodePort,omitempty"`
	PassengerMaxPoolSize    int      `json:"passengerMaxPoolSize,omitempty"`
	InitialDelaySeconds     int      `json:"initialDelaySeconds,omitempty"`
	ReadinessPeriodSeconds  int      `json:"readinessPeriodSeconds,omitempty"`
	ReadinessTimeoutSeconds int      `json:"readinessTimeoutSeconds,omitempty"`
	FailureThreshold        int      `json:"failureThreshold,omitempty"`
	Hpa                     Hpa      `json:"hpa,omitempty"`
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

type Nomex struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
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
	MetagpuEnabled       bool              `json:"metagpuEnabled,omitempty"`
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
	Nomex                           Nomex                           `json:"nomex,omitempty"`
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

type PriorityClass struct {
	Name        string `json:"name"`
	Value       int32  `json:"value"`
	Description string `json:"description"`
}

type Tenancy struct {
	Enabled bool   `json:"enabled,omitempty"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
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

type HugePages struct {
	Enabled bool   `json:"enabled,omitempty"`
	Size    string `json:"size,omitempty"`
	Memory  string `json:"memory,omitempty"`
}

type ExtraScrapeConfigs struct {
	Role          string `json:"role,omitempty"`
	Namespace     string `json:"namespace,omitempty"`
	LabelSelector string `json:"labelSelector,omitempty"`
}

type Grafana struct {
	Enabled  bool   `json:"enabled,omitempty"`
	Image    string `json:"image,omitempty"`
	SvcName  string `json:"svcName,omitempty"`
	Port     int    `json:"port,omitempty"`
	NodePort int    `json:"nodePort,omitempty"`
	CredsRef string `json:"credsRef,omitempty"`
}

type Prom struct {
	Enabled            bool                  `json:"enabled,omitempty"`
	SvcName            string                `json:"svcName,omitempty"`
	CredsRef           string                `json:"credsRef,omitempty"`
	ExtraScrapeConfigs []*ExtraScrapeConfigs `json:"extraScrapeConfigs,omitempty"`
	Image              string                `json:"image,omitempty"`
	Grafana            Grafana               `json:"grafana,omitempty"`
}

type Pg struct {
	Enabled            bool              `json:"enabled,omitempty"`
	ServiceAccount     string            `json:"serviceAccount,omitempty"`
	Image              string            `json:"image,omitempty"`
	Port               int               `json:"port,omitempty"`
	StorageSize        string            `json:"storageSize,omitempty"`
	SvcName            string            `json:"svcName,omitempty"`
	StorageClass       string            `json:"storageClass,omitempty"`
	Requests           Requests          `json:"requests,omitempty"`
	Limits             Limits            `json:"limits,omitempty"`
	MaxConnections     int               `json:"maxConnections,omitempty"`
	SharedBuffers      string            `json:"sharedBuffers,omitempty"`      // https://github.com/sclorg/postgresql-container/tree/generated/12
	EffectiveCacheSize string            `json:"effectiveCacheSize,omitempty"` // https://github.com/sclorg/postgresql-container/tree/generated/12
	HugePages          HugePages         `json:"hugePages,omitempty"`
	NodeSelector       map[string]string `json:"nodeSelector,omitempty"`
	CredsRef           string            `json:"credsRef,omitempty"`
	PvcName            string            `json:"pvcName,omitempty"`
}

type Minio struct {
	Enabled        bool              `json:"enabled,omitempty"`
	ServiceAccount string            `json:"serviceAccount,omitempty"`
	Replicas       int               `json:"replicas,omitempty"`
	Image          string            `json:"image,omitempty"`
	Port           int               `json:"port,omitempty"`
	StorageSize    string            `json:"storageSize,omitempty"`
	SvcName        string            `json:"svcName,omitempty"`
	NodePort       int               `json:"nodePort,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	Requests       Requests          `json:"requests,omitempty"`
	Limits         Limits            `json:"limits,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	PvcName        string            `json:"pvcName,omitempty"`
}

type Redis struct {
	Enabled        bool              `json:"enabled,omitempty"`
	ServiceAccount string            `json:"serviceAccount,omitempty"`
	Image          string            `json:"image,omitempty"`
	SvcName        string            `json:"svcName,omitempty"`
	Port           int               `json:"port,omitempty"`
	StorageSize    string            `json:"storageSize,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	Requests       Requests          `json:"requests,omitempty"`
	Limits         Limits            `json:"limits,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	CredsRef       string            `json:"credsRef,omitempty"`
	PvcName        string            `json:"pvcName,omitempty"`
}

type Kibana struct {
	Enabled        bool     `json:"enabled,omitempty"`
	ServiceAccount string   `json:"serviceAccount,omitempty"`
	SvcName        string   `json:"svcName,omitempty"`
	Port           int      `json:"port,omitempty"`
	Image          string   `json:"image,omitempty"`
	NodePort       int      `json:"nodePort,omitempty"`
	Requests       Requests `json:"requests,omitempty"`
	Limits         Limits   `json:"limits,omitempty"`
	CredsRef       string   `json:"credsRef,omitempty"`
}

type Es struct {
	Enabled        bool              `json:"enabled,omitempty"`
	ServiceAccount string            `json:"serviceAccount,omitempty"`
	Image          string            `json:"image,omitempty"`
	Port           int               `json:"port,omitempty"`
	StorageSize    string            `json:"storageSize,omitempty"`
	SvcName        string            `json:"svcName,omitempty"`
	NodePort       int               `json:"nodePort,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	Requests       Requests          `json:"requests,omitempty"`
	Limits         Limits            `json:"limits,omitempty"`
	JavaOpts       string            `json:"javaOpts,omitempty"`
	PatchEsNodes   bool              `json:"patchEsNodes,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	CredsRef       string            `json:"credsRef,omitempty"`
	PvcName        string            `json:"pvcName,omitempty"`
	CleanupPolicy  CleanupPolicy     `json:"cleanupPolicy,omitempty"`
	Kibana         Kibana            `json:"kibana,omitempty"`
	Elastalert     Elastalert        `json:"elastalert,omitempty"`
}

type CleanupPolicy struct {
	All       string `json:"all,omitempty"`
	App       string `json:"app,omitempty"`
	Jobs      string `json:"jobs,omitempty"`
	Endpoints string `json:"endpoints,omitempty"`
}

type Dbs struct {
	Pg    Pg    `json:"pg,omitempty"`
	Redis Redis `json:"redis,omitempty"`
	Minio Minio `json:"minio,omitempty"`
	Es    Es    `json:"es,omitempty"`
	Cvat  Cvat  `json:"cvat,omitempty"`
	Prom  Prom  `json:"prom,omitempty"`
}

type Cvat struct {
	Enabled bool  `json:"enabled,omitempty"`
	Pg      Pg    `json:"pg,omitempty"`
	Redis   Redis `json:"redis,omitempty"`
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
	CertSecret string `json:"certSecret,omitempty"`
}

type Proxy struct {
	Enabled    bool     `json:"enabled,omitempty"`
	ConfigRef  string   `json:"configRef,omitempty"`
	HttpProxy  []string `json:"httpProxy,omitempty"`
	HttpsProxy []string `json:"httpsProxy,omitempty"`
	NoProxy    []string `json:"noProxy,omitempty"`
}

type Networking struct {
	Ingress Ingress `json:"ingress,omitempty"`
	HTTPS   HTTPS   `json:"https,omitempty"`
	Proxy   Proxy   `json:"proxy,omitempty"`
}

type SSO struct {
	Enabled bool       `json:"enabled,omitempty"`
	Pki     Pki        `json:"pki,omitempty"`
	Jwks    Jwks       `json:"jwks,omitempty"`
	Central CentralSSO `json:"central,omitempty"`
	Authz   Authz      `json:"authz,omitempty"`
}

type Pki struct {
	Enabled          bool   `json:"enabled,omitempty"`
	RootCaSecret     string `json:"rootCaSecret,omitempty"`
	PrivateKeySecret string `json:"privateKeySecret,omitempty"`
	PublicKeySecret  string `json:"publicKeySecret,omitempty"`
}

type JwksCache struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type Jwks struct {
	Enabled bool      `json:"enabled,omitempty"`
	Name    string    `json:"name,omitempty"`
	Image   string    `json:"image,omitempty"`
	Cache   JwksCache `json:"cache,omitempty"`
}

type Authz struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
	Address string `json:"address,omitempty"`
}

type CentralSSO struct {
	Enabled                          bool     `json:"enabled,omitempty"`
	PublicUrl                        string   `json:"publicUrl,omitempty"`
	CnvrgProxyImage                  string   `json:"cnvrgProxyImage,omitempty"`
	OauthProxyImage                  string   `json:"oauthProxyImage,omitempty"`
	CentralUiImage                   string   `json:"centralUiImage,omitempty"`
	AdminUser                        string   `json:"adminUser,omitempty"`
	Provider                         string   `json:"provider,omitempty"`
	EmailDomain                      []string `json:"emailDomain,omitempty"`
	ClientID                         string   `json:"clientId,omitempty"`
	ClientSecret                     string   `json:"clientSecret,omitempty"`
	OidcIssuerURL                    string   `json:"oidcIssuerUrl,omitempty"`
	ServiceUrl                       string   `json:"serviceUrl,omitempty"`
	Scope                            string   `json:"scope,omitempty"`
	InsecureOidcAllowUnverifiedEmail bool     `json:"insecureOidcAllowUnverifiedEmail,omitempty"`
	WhitelistDomain                  string   `json:"whitelistDomain,omitempty"`
	CookieDomain                     string   `json:"cookieDomain,omitempty"`
	GroupsAuth                       bool     `json:"groupsAuth,omitempty"`
}

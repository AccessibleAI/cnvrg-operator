package v1

type HugePages struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Size    string `json:"size,omitempty"`
	Memory  string `json:"memory,omitempty"`
}

type Pg struct {
	Enabled        *bool             `json:"enabled,omitempty"`
	ServiceAccount string            `json:"serviceAccount,omitempty"`
	Image          string            `json:"image,omitempty"`
	Port           int               `json:"port,omitempty"`
	StorageSize    string            `json:"storageSize,omitempty"`
	SvcName        string            `json:"svcName,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	Requests       Requests          `json:"requests,omitempty"`
	MaxConnections int               `json:"maxConnections,omitempty"`
	SharedBuffers  string            `json:"sharedBuffers,omitempty"`
	HugePages      HugePages         `json:"hugePages,omitempty"`
	Fixpg          *bool             `json:"fixpg,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	CredsRef       string            `json:"credsRef,omitempty"`
	PvcName        string            `json:"pvcName,omitempty"`
}

type Minio struct {
	Enabled        *bool             `json:"enabled,omitempty"`
	ServiceAccount string            `json:"serviceAccount,omitempty"`
	Replicas       int               `json:"replicas,omitempty"`
	Image          string            `json:"image,omitempty"`
	Port           int               `json:"port,omitempty"`
	StorageSize    string            `json:"storageSize,omitempty"`
	SvcName        string            `json:"svcName,omitempty"`
	NodePort       int               `json:"nodePort,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	Requests       Requests          `json:"requests,omitempty"`
	SharedStorage  SharedStorage     `json:"sharedStorage,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	PvcName        string            `json:"pvcName,omitempty"`
}

type Redis struct {
	Enabled        *bool             `json:"enabled,omitempty"`
	ServiceAccount string            `json:"serviceAccount,omitempty"`
	Image          string            `json:"image,omitempty"`
	SvcName        string            `json:"svcName,omitempty"`
	Port           int               `json:"port,omitempty"`
	StorageSize    string            `json:"storageSize,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	Limits         Limits            `json:"limits,omitempty"`
	Requests       Requests          `json:"requests,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	CredsRef       string            `json:"credsRef,omitempty"`
	PvcName        string            `json:"pvcName,omitempty"`
}

type Es struct {
	Enabled        *bool             `json:"enabled,omitempty"`
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
	PatchEsNodes   *bool             `json:"patchEsNodes,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	CredsRef       string            `json:"credsRef,omitempty"`
	PvcName        string            `json:"pvcName,omitempty"`
}

type AppDbs struct {
	Pg    Pg    `json:"pg,omitempty"`
	Redis Redis `json:"redis,omitempty"`
	Minio Minio `json:"minio,omitempty"`
	Es    Es    `json:"es,omitempty"`
}

type InfraDbs struct {
	Redis Redis `json:"redis,omitempty"`
}

var minioDefaults = Minio{
	Enabled:        &defaultEnabled,
	ServiceAccount: "minio",
	Replicas:       1,
	Image:          "docker.io/minio/minio:RELEASE.2020-09-17T04-49-20Z",
	Port:           9000,
	StorageSize:    "100Gi",
	SvcName:        "minio",
	NodePort:       30090,
	StorageClass:   "",
	Requests: Requests{
		Cpu:    "1000m",
		Memory: "2Gi",
	},
	PvcName: "",
	SharedStorage: SharedStorage{
		Enabled:          &defaultEnabled,
		UseExistingClaim: "",
		ConsistentHash: ConsistentHash{
			Key:   "httpQueryParameterName",
			Value: "uploadId",
		},
	},
}

var pgDefault = Pg{
	Enabled:        &defaultEnabled,
	ServiceAccount: "pg",
	Image:          "centos/postgresql-12-centos7",
	Port:           5432,
	StorageSize:    "80Gi",
	SvcName:        "postgres",
	StorageClass:   "",
	Requests: Requests{
		Cpu:    "4000m",
		Memory: "4Gi",
	},
	MaxConnections: 500,
	SharedBuffers:  "64MB",
	Fixpg:          &defaultTrue,
	NodeSelector:   nil,
	PvcName:        "",
	HugePages: HugePages{
		Enabled: &defaultEnabled,
		Size:    "2Mi",
		Memory:  "",
	},
	CredsRef: "pg-creds",
}

var redisDefault = Redis{
	Enabled:        &defaultEnabled,
	ServiceAccount: "redis",
	Image:          "docker.io/cnvrg/cnvrg-redis:v3.0.5.c2",
	SvcName:        "redis",
	Port:           6379,
	StorageSize:    "10Gi",
	StorageClass:   "",
	NodeSelector:   nil,
	CredsRef:       "redis-creds",
	PvcName:        "",
	Limits: Limits{
		Cpu:    "1000m",
		Memory: "2Gi",
	},
	Requests: Requests{
		Cpu:    "100m",
		Memory: "200Mi",
	},
}

var esDefault = Es{
	Enabled:        &defaultEnabled,
	ServiceAccount: "es",
	Image:          "docker.io/cnvrg/cnvrg-es:v7.8.1.a1",
	Port:           9200,
	StorageSize:    "30Gi",
	SvcName:        "elasticsearch",
	NodePort:       32200,
	StorageClass:   "",
	Requests: Requests{
		Cpu:    "1000m",
		Memory: "1Gi",
	},
	Limits: Limits{
		Cpu:    "2000m",
		Memory: "4Gi",
	},
	JavaOpts:     "",
	PatchEsNodes: &defaultTrue,
	CredsRef:     "es-creds",
}

var appDbsDefaults = AppDbs{
	Pg:    pgDefault,
	Redis: redisDefault,
	Minio: minioDefaults,
	Es:    esDefault,
}

var infraDbsDefaults = InfraDbs{
	Redis: redisDefault,
}

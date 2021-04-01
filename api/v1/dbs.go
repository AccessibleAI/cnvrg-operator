package v1

type HugePages struct {
	Enabled string `json:"enabled,omitempty"`
	Size    string `json:"size,omitempty"`
	Memory  string `json:"memory,omitempty"`
}

type Pg struct {
	Enabled        string            `json:"enabled,omitempty"`
	ServiceAccount string            `json:"serviceAccount"`
	SecretName     string            `json:"secretName,omitempty"`
	Image          string            `json:"image,omitempty"`
	Port           int               `json:"port,omitempty"`
	StorageSize    string            `json:"storageSize,omitempty"`
	SvcName        string            `json:"svcName,omitempty"`
	Dbname         string            `json:"dbname,omitempty"`
	Pass           string            `json:"pass,omitempty"`
	User           string            `json:"user,omitempty"`
	RunAsUser      int               `json:"runAsUser,omitempty"`
	FsGroup        int               `json:"fsGroup,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	CPURequest     int               `json:"cpuRequest,omitempty"`
	MemoryRequest  string            `json:"memoryRequest,omitempty"`
	MaxConnections int               `json:"maxConnections,omitempty"`
	SharedBuffers  string            `json:"sharedBuffers,omitempty"`
	HugePages      HugePages         `json:"hugePages,omitempty"`
	Fixpg          string            `json:"fixpg,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	Tolerations    map[string]string `json:"tolerations,omitempty"`
}

type Minio struct {
	Enabled        string            `json:"enabled,omitempty"`
	ServiceAccount string            `json:"serviceAccount,omitempty"`
	Replicas       int               `json:"replicas,omitempty"`
	Image          string            `json:"image,omitempty"`
	Port           int               `json:"port,omitempty"`
	StorageSize    string            `json:"storageSize,omitempty"`
	SvcName        string            `json:"svcName,omitempty"`
	NodePort       int               `json:"nodePort,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	CPURequest     int               `json:"cpuRequest,omitempty"`
	MemoryRequest  string            `json:"memoryRequest,omitempty"`
	SharedStorage  SharedStorage     `json:"sharedStorage,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	Tolerations    map[string]string `json:"tolerations,omitempty"`
}

type Redis struct {
	Enabled        string            `json:"enabled,omitempty"`
	ServiceAccount string            `json:"serviceAccount"`
	Image          string            `json:"image,omitempty"`
	SvcName        string            `json:"svcName,omitempty"`
	Port           int               `json:"port,omitempty"`
	Appendonly     string            `json:"appendonly,omitempty"`
	StorageSize    string            `json:"storageSize,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	Limits         Limits            `json:"limits,omitempty"`
	Requests       Requests          `json:"requests,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	Tolerations    map[string]string `json:"tolerations,omitempty"`
}

type Es struct {
	Enabled        string            `json:"enabled,omitempty"`
	ServiceAccount string            `json:"serviceAccount,omitempty"`
	Image          string            `json:"image,omitempty"`
	Port           int               `json:"port,omitempty"`
	StorageSize    string            `json:"storageSize,omitempty"`
	SvcName        string            `json:"svcName,omitempty"`
	RunAsUser      int               `json:"runAsUser,omitempty"`
	FsGroup        int               `json:"fsGroup,omitempty"`
	NodePort       int               `json:"nodePort,omitempty"`
	StorageClass   string            `json:"storageClass,omitempty"`
	CPURequest     int               `json:"cpuRequest,omitempty"`
	MemoryRequest  string            `json:"memoryRequest,omitempty"`
	CPULimit       int               `json:"cpuLimit,omitempty"`
	MemoryLimit    string            `json:"memoryLimit,omitempty"`
	JavaOpts       string            `json:"javaOpts,omitempty"`
	PatchEsNodes   string            `json:"patchEsNodes,omitempty"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty"`
	Tolerations    map[string]string `json:"tolerations,omitempty"`
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
	Enabled:       "true",
	ServiceAccount: "default",
	Replicas:      1,
	Image:         "docker.io/minio/minio:RELEASE.2020-09-17T04-49-20Z",
	Port:          9000,
	StorageSize:   "100Gi",
	SvcName:       "minio",
	NodePort:      30090,
	StorageClass:  "",
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

var pgDefault = Pg{
	Enabled:        "true",
	ServiceAccount: "default",
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
	StorageClass:   "",
	CPURequest:     4,
	MemoryRequest:  "4Gi",
	MaxConnections: 100,
	SharedBuffers:  "64MB",
	Fixpg:          "true",
	NodeSelector:   nil,
	Tolerations:    nil,
	HugePages: HugePages{
		Enabled: "false",
		Size:    "2Mi",
		Memory:  "",
	},
}

var redisDefault = Redis{
	Enabled:        "true",
	ServiceAccount: "default",
	Image:          "docker.io/cnvrg/cnvrg-redis:v3.0.5.c2",
	SvcName:        "redis",
	Port:           6379,
	Appendonly:     "yes",
	StorageSize:    "10Gi",
	StorageClass:   "",
	NodeSelector:   nil,
	Tolerations:    nil,
	Limits: Limits{
		CPU:    1,
		Memory: "2Gi",
	},
	Requests: Requests{
		CPU:    "100m",
		Memory: "200Mi",
	},
}

var esDefault = Es{
	Enabled:        "true",
	ServiceAccount: "es",
	Image:          "docker.io/cnvrg/cnvrg-es:v7.8.1",
	Port:           9200,
	StorageSize:    "30Gi",
	SvcName:        "elasticsearch",
	RunAsUser:      1000,
	FsGroup:        1000,
	NodePort:       32200,
	StorageClass:   "",
	CPURequest:     1,
	MemoryRequest:  "1Gi",
	CPULimit:       2,
	MemoryLimit:    "4Gi",
	JavaOpts:       "",
	PatchEsNodes:   "true",
	NodeSelector:   nil,
	Tolerations:    nil,
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

package v1

type HugePages struct {
	Enabled bool   `json:"enabled,omitempty"`
	Size    string `json:"size,omitempty"`
	Memory  string `json:"memory,omitempty"`
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
	Backup             Backup            `json:"backup,omitempty"`
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
	SharedStorage  SharedStorage     `json:"sharedStorage,omitempty"`
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
	Backup         Backup            `json:"backup,omitempty"`
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
}

type CleanupPolicy struct {
	All       string `json:"all,omitempty"`
	App       string `json:"app,omitempty"`
	Jobs      string `json:"jobs,omitempty"`
	Endpoints string `json:"endpoints,omitempty"`
}

type AppDbs struct {
	Pg    Pg    `json:"pg,omitempty"`
	Redis Redis `json:"redis,omitempty"`
	Minio Minio `json:"minio,omitempty"`
	Es    Es    `json:"es,omitempty"`
	Cvat  Cvat  `json:"cvat,omitempty"`
}

type Cvat struct {
	Enabled bool  `json:"enabled,omitempty"`
	Pg      Pg    `json:"pg,omitempty"`
	Redis   Redis `json:"redis,omitempty"`
}

type InfraDbs struct {
	Redis Redis `json:"redis,omitempty"`
}

var minioDefaults = Minio{
	Enabled:        false,
	ServiceAccount: "minio",
	Replicas:       1,
	Image:          "minio:RELEASE.2021-05-22T02-34-39Z",
	Port:           9000,
	StorageSize:    "100Gi",
	SvcName:        "minio",
	NodePort:       30090,
	StorageClass:   "",
	Requests: Requests{
		Cpu:    "200m",
		Memory: "2Gi",
	},
	Limits: Limits{
		Cpu:    "8",
		Memory: "20Gi",
	},
	PvcName: "minio-storage",
	SharedStorage: SharedStorage{
		Enabled: false,
		ConsistentHash: ConsistentHash{
			Key:   "httpQueryParameterName",
			Value: "uploadId",
		},
	},
}

var pgDefault = Pg{
	Enabled:        false,
	ServiceAccount: "pg",
	Image:          "postgresql-12-centos7:latest",
	Port:           5432,
	StorageSize:    "80Gi",
	SvcName:        "postgres",
	StorageClass:   "",
	Requests: Requests{
		Cpu:    "1",
		Memory: "4Gi",
	},
	Limits: Limits{
		Cpu:    "12",
		Memory: "32Gi",
	},
	MaxConnections:     500,
	SharedBuffers:      "1024MB", // for the shared_buffers we use 1/4 of given memory
	EffectiveCacheSize: "2048MB", // for the effective_cache_size we set the value to 1/2 of the given memory
	NodeSelector:       nil,
	PvcName:            "pg-storage",
	HugePages: HugePages{
		Enabled: false,
		Size:    "2Mi", // 2Mi/1Gi https://kubernetes.io/docs/tasks/manage-hugepages/scheduling-hugepages/ ,  https://wiki.debian.org/Hugepages#Huge_pages_sizes
		Memory:  "",
	},
	CredsRef: "pg-creds",
	Backup: Backup{
		Enabled:   false,
		BucketRef: "cp-object-storage",
		CredsRef:  "pg-creds",
		Rotation:  5,
		Period:    "24h",
	},
}

var redisDefault = Redis{
	Enabled:        false,
	ServiceAccount: "redis",
	Image:          "cnvrg-redis:v3.0.5.c2",
	SvcName:        "redis",
	Port:           6379,
	StorageSize:    "10Gi",
	StorageClass:   "",
	NodeSelector:   nil,
	CredsRef:       "redis-creds",
	PvcName:        "redis-storage",
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
	Enabled:        false,
	ServiceAccount: "es",
	Image:          "cnvrg-es:v7.8.1.a1-dynamic-indices",
	Port:           9200,
	StorageSize:    "80Gi",
	SvcName:        "elasticsearch",
	NodePort:       32200,
	StorageClass:   "",
	Requests: Requests{
		Cpu:    "500m",
		Memory: "4Gi",
	},
	Limits: Limits{
		Cpu:    "4",
		Memory: "8Gi",
	},
	JavaOpts:     "",
	PatchEsNodes: false,
	CredsRef:     "es-creds",
	PvcName:      "es-storage",
	CleanupPolicy: CleanupPolicy{
		All: "3d",
		App: "30d",
		Jobs: "14d",
		Endpoints: "1825d",
	},
}

var appDbsDefaults = AppDbs{
	Pg:    pgDefault,
	Redis: redisDefault,
	Minio: minioDefaults,
	Es:    esDefault,
	Cvat:  cvatDefault,
}

var cvatDefault = Cvat{
	Enabled: false,
	Pg:      cvatPgDefault,
	Redis:   cvatRedisDefault,
}

var cvatPgDefault = Pg{
	Enabled:        false,
	ServiceAccount: "cvat-pg",
	Image:          "postgresql-12-centos7:latest",
	Port:           5432,
	StorageSize:    "100Gi",
	SvcName:        "cvat-postgres",
	StorageClass:   "",
	Requests: Requests{
		Cpu:    "1",
		Memory: "2Gi",
	},
	Limits: Limits{
		Cpu:    "2",
		Memory: "4Gi",
	},
	MaxConnections:     500,
	SharedBuffers:      "1024MB", // for the shared_buffers we use 1/4 of given memory
	EffectiveCacheSize: "2048MB", // for the effective_cache_size we set the value to 1/2 of the given memory
	NodeSelector:       nil,
	PvcName:            "cvat-pg-storage",
	HugePages: HugePages{
		Enabled: false,
		Size:    "2Mi", // 2Mi/1Gi https://kubernetes.io/docs/tasks/manage-hugepages/scheduling-hugepages/ ,  https://wiki.debian.org/Hugepages#Huge_pages_sizes
		Memory:  "",
	},
	CredsRef: "cvat-pg-secret",
	Backup: Backup{
		Enabled:   false,
		BucketRef: "cp-object-storage",
		CredsRef:  "cvat-pg-secret",
		Rotation:  5,
		Period:    "24h",
	},
}

var cvatRedisDefault = Redis{
	Enabled:        false,
	ServiceAccount: "cvat-redis",
	Image:          "redis:4.0.5-alpine",
	SvcName:        "cvat-redis",
	Port:           6379,
	StorageSize:    "10Gi",
	StorageClass:   "",
	NodeSelector:   nil,
	CredsRef:       "cvat-redis-secret",
	PvcName:        "cvat-redis-storage",
	Limits: Limits{
		Cpu:    "1000m",
		Memory: "2Gi",
	},
	Requests: Requests{
		Cpu:    "250m",
		Memory: "500Mi",
	},
}

var infraDbsDefaults = InfraDbs{
	Redis: redisDefault,
}

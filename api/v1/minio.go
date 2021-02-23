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

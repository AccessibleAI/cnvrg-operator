package v1

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

package v1

type Backup struct {
	Enabled   bool   `json:"enabled,omitempty"`
	BucketRef string `json:"bucketRef,omitempty"`
	CredsRef  string `json:"credsRef,omitempty"`
	Rotation  int    `json:"rotation,omitempty"`
	Period    string `json:"period,omitempty"` // on of [Xs, Xm, Xh]
}

type Capsule struct {
	Enabled      bool     `json:"enabled,omitempty"`
	Image        string   `json:"image,omitempty"`
	Requests     Requests `json:"requests,omitempty"`
	Limits       Limits   `json:"limits,omitempty"`
	SvcName      string   `json:"svcName,omitempty"`
	StorageSize  string   `json:"storageSize,omitempty"`
	StorageClass string   `json:"storageClass,omitempty"`
}

var capsuleDefault = Capsule{
	Enabled: false,
	Image:   "cnvrg-capsule:1.0.2",
	Requests: Requests{
		Cpu:    "200m",
		Memory: "500Mi",
	},
	Limits: Limits{
		Cpu:    "1",
		Memory: "1Gi",
	},
	SvcName:      "capsule",
	StorageSize:  "100Gi",
	StorageClass: "",
}

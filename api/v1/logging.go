package v1

type Fluentbit struct {
	Image string `json:"image,omitempty"`
}

type Es struct {
	Enabled       string `json:"enabled,omitempty"`
	Image         string `json:"image,omitempty"`
	Port          int    `json:"port,omitempty"`
	StorageSize   string `json:"storageSize,omitempty"`
	SvcName       string `json:"svcName,omitempty"`
	RunAsUser     int    `json:"runAsUser,omitempty"`
	FsGroup       int    `json:"fsGroup,omitempty"`
	NodePort      int    `json:"nodePort,omitempty"`
	StorageClass  string `json:"storageClass,omitempty"`
	CPURequest    int    `json:"cpuRequest,omitempty"`
	MemoryRequest string `json:"memoryRequest,omitempty"`
	CPULimit      int    `json:"cpuLimit,omitempty"`
	MemoryLimit   string `json:"memoryLimit,omitempty"`
	JavaOpts      string `json:"javaOpts,omitempty"`
	PatchEsNodes  string `json:"patchEsNodes,omitempty"`
}

type Elastalert struct {
	Enabled       string `json:"enabled,omitempty"`
	Image         string `json:"image,omitempty"`
	Port          int    `json:"port,omitempty"`
	NodePort      int    `json:"nodePort,omitempty"`
	ContainerPort int    `json:"containerPort,omitempty"`
	StorageSize   string `json:"storageSize,omitempty"`
	SvcName       string `json:"svcName,omitempty"`
	StorageClass  string `json:"storageClass,omitempty"`
	CPURequest    string `json:"cpuRequest,omitempty"`
	MemoryRequest string `json:"memoryRequest,omitempty"`
	CPULimit      string `json:"cpuLimit,omitempty"`
	MemoryLimit   string `json:"memoryLimit,omitempty"`
	RunAsUser     int    `json:"runAsUser,omitempty"`
	FsGroup       int    `json:"fsGroup,omitempty"`
}

type Kibana struct {
	Enabled       string `json:"enabled,omitempty"`
	SvcName       string `json:"svcName,omitempty"`
	Port          int    `json:"port,omitempty"`
	Image         string `json:"image,omitempty"`
	NodePort      int    `json:"nodePort,omitempty"`
	CPURequest    string `json:"cpuRequest,omitempty"`
	MemoryRequest string `json:"memoryRequest,omitempty"`
	CPULimit      int    `json:"cpuLimit,omitempty"`
	MemoryLimit   string `json:"memoryLimit,omitempty"`
}

type CnvrgAppLogging struct {
	Enabled    string     `json:"enabled,omitempty"`
	Es         Es         `json:"es,omitempty"`
	Elastalert Elastalert `json:"elastalert,omitempty"`
	Kibana     Kibana     `json:"kibana,omitempty"`
}

type CnvrgInfraLogging struct {
	Enabled   string    `json:"enabled,omitempty"`
	Fluentbit Fluentbit `json:"fluentbit,omitempty"`
}

var fluentbitDefault = Fluentbit{
	Image: "cnvrg/cnvrg-fluentbit:v1.7.2",
}

var cnvrgAppLoggingDefault = CnvrgAppLogging{
	Enabled: "true",
	Es: Es{
		Enabled:       "true",
		Image:         "docker.io/cnvrg/cnvrg-es:v7.8.1",
		Port:          9200,
		StorageSize:   "30Gi",
		SvcName:       "elasticsearch",
		RunAsUser:     1000,
		FsGroup:       1000,
		NodePort:      32200,
		StorageClass:  "use-default",
		CPURequest:    1,
		MemoryRequest: "1Gi",
		CPULimit:      2,
		MemoryLimit:   "4Gi",
		JavaOpts:      "",
		PatchEsNodes:  "true",
	},
	Elastalert: Elastalert{
		Enabled:       "true",
		Image:         "bitsensor/elastalert:3.0.0-beta.1",
		Port:          80,
		NodePort:      32030,
		ContainerPort: 3030,
		StorageSize:   "30Gi",
		SvcName:       "elastalert",
		StorageClass:  "use-default",
		CPURequest:    "100m",
		MemoryRequest: "200Mi",
		CPULimit:      "400m",
		MemoryLimit:   "800Mi",
		RunAsUser:     1000,
		FsGroup:       1000,
	},
	Kibana: Kibana{
		Enabled:       "true",
		SvcName:       "kibana",
		Port:          5601,
		Image:         "docker.elastic.co/kibana/kibana-oss:7.8.1",
		NodePort:      30601,
		CPURequest:    "100m",
		MemoryRequest: "100Mi",
		CPULimit:      1,
		MemoryLimit:   "2Gi",
	},
}

var cnvrgInfraLoggingDefault = CnvrgInfraLogging{
	Enabled:   "true",
	Fluentbit: fluentbitDefault,
}

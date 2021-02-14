package pg

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
}

package v1

type Storage struct {
	Enabled         string   `json:"enabled,omitempty"`
	Hostpath        Hostpath `json:"hostpath,omitempty"`
	Nfs             Nfs      `json:"nfs,omitempty"`
}

type Hostpath struct {
	Enabled          string `json:"enabled,omitempty"`
	Image            string `json:"image,omitempty"`
	HostPath         string `json:"hostPath,omitempty"`
	StorageClassName string `json:"storageClassName,omitempty"`
	NodeName         string `json:"nodeName,omitempty"`
	CPURequest       string `json:"cpuRequest,omitempty"`
	MemoryRequest    string `json:"memoryRequest,omitempty"`
	CPULimit         string `json:"cpuLimit,omitempty"`
	MemoryLimit      string `json:"memoryLimit,omitempty"`
	ReclaimPolicy    string `json:"reclaimPolicy,omitempty"`
	DefaultSc        string `json:"defaultSc,omitempty"`
}

type Nfs struct {
	Enabled          string `json:"enabled,omitempty"`
	Image            string `json:"image,omitempty"`
	Provisioner      string `json:"provisioner,omitempty"`
	StorageClassName string `json:"storageClassName,omitempty"`
	Server           string `json:"server,omitempty"`
	Path             string `json:"path,omitempty"`
	CPURequest       string `json:"cpuRequest,omitempty"`
	MemoryRequest    string `json:"memoryRequest,omitempty"`
	CPULimit         string `json:"cpuLimit,omitempty"`
	MemoryLimit      string `json:"memoryLimit,omitempty"`
	ReclaimPolicy    string `json:"reclaimPolicy,omitempty"`
	DefaultSc        string `json:"defaultSc,omitempty"`
}

var storageDefault = Storage{
	Enabled:         "false",
	Hostpath: Hostpath{
		Enabled:          "false",
		Image:            "quay.io/kubevirt/hostpath-provisioner",
		HostPath:         "/cnvrg-hostpath-storage",
		StorageClassName: "cnvrg-hostpath-storage",
		NodeName:         "",
		CPURequest:       "100m",
		MemoryRequest:    "100Mi",
		CPULimit:         "200m",
		MemoryLimit:      "200Mi",
		ReclaimPolicy:    "Retain",
		DefaultSc:        "false",
	},
	Nfs: Nfs{
		Enabled:          "false",
		Image:            "gcr.io/k8s-staging-sig-storage/nfs-subdir-external-provisioner:v4.0.0",
		Provisioner:      "cnvrg.io/ifs",
		StorageClassName: "cnvrg-nfs-storage",
		Server:           "",
		Path:             "",
		CPURequest:       "100m",
		MemoryRequest:    "100Mi",
		CPULimit:         "100m",
		MemoryLimit:      "200Mi",
		ReclaimPolicy:    "Retain",
		DefaultSc:        "false",
	},
}

package v1

type Storage struct {
	Hostpath Hostpath `json:"hostpath,omitempty"`
	Nfs      Nfs      `json:"nfs,omitempty"`
}

type Hostpath struct {
	Enabled          bool              `json:"enabled,omitempty"`
	Image            string            `json:"image,omitempty"`
	Path             string            `json:"path,omitempty"`
	StorageClassName string            `json:"storageClassName,omitempty"`
	Requests         Requests          `json:"requests,omitempty"`
	Limits           Limits            `json:"limits,omitempty"`
	ReclaimPolicy    string            `json:"reclaimPolicy,omitempty"`
	DefaultSc        bool              `json:"defaultSc,omitempty"`
	NodeSelector     map[string]string `json:"nodeSelector,omitempty"`
}

type Nfs struct {
	Enabled          bool     `json:"enabled,omitempty"`
	Image            string   `json:"image,omitempty"`
	Provisioner      string   `json:"provisioner,omitempty"`
	StorageClassName string   `json:"storageClassName,omitempty"`
	Server           string   `json:"server,omitempty"`
	Path             string   `json:"path,omitempty"`
	Requests         Requests `json:"requests,omitempty"`
	Limits           Limits   `json:"limits,omitempty"`
	ReclaimPolicy    string   `json:"reclaimPolicy,omitempty"`
	DefaultSc        bool     `json:"defaultSc,omitempty"`
}

var storageDefault = Storage{
	Hostpath: Hostpath{
		Enabled:          false,
		Image:            "hostpath-provisioner:latest",
		Path:             "/cnvrg-hostpath-storage",
		StorageClassName: "cnvrg-hostpath-storage",
		Requests: Requests{
			Cpu:    "100m",
			Memory: "100Mi",
		},
		Limits: Limits{
			Cpu:    "200m",
			Memory: "200Mi",
		},
		ReclaimPolicy: "Retain",
		DefaultSc:     false,
		NodeSelector:  nil,
	},
	Nfs: Nfs{
		Enabled:          false,
		Image:            "nfs-subdir-external-provisioner:v4.0.0",
		Provisioner:      "cnvrg.io/ifs",
		StorageClassName: "cnvrg-nfs-storage",
		Server:           "",
		Path:             "",
		Requests: Requests{
			Cpu:    "100m",
			Memory: "100Mi",
		},
		Limits: Limits{
			Cpu:    "200m",
			Memory: "200Mi",
		},
		ReclaimPolicy: "Retain",
		DefaultSc:     false,
	},
}

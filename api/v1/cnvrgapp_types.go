/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

type Storage struct {
	Enabled         string   `json:"enabled,omitempty"`
	CcpStorageClass string   `json:"ccpStorageClass,omitempty"`
	Hostpath        Hostpath `json:"hostpath,omitempty"`
	Nfs             Nfs      `json:"nfs,omitempty"`
}

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CnvrgAppSpec defines the desired state of CnvrgApp
type CnvrgAppSpec struct {
	Pg      Pg      `json:"pg,omitempty"`
	Storage Storage `json:"storage,omitempty"`
	Message string  `json:"message,omitempty"`
}

// CnvrgAppStatus defines the observed state of CnvrgApp
type CnvrgAppStatus struct {
	Message string `json:"message,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// CnvrgApp is the Schema for the cnvrgapps API
type CnvrgApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CnvrgAppSpec   `json:"spec,omitempty"`
	Status CnvrgAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CnvrgAppList contains a list of CnvrgApp
type CnvrgAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CnvrgApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CnvrgApp{}, &CnvrgAppList{})
}

func DefaultCnvrgAppSpec() CnvrgAppSpec {
	return CnvrgAppSpec{
		Pg: Pg{
			Enabled:        "true",
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
			StorageClass:   "use-default",
			CPURequest:     4,
			MemoryRequest:  "4Gi",
			MaxConnections: 100,
			SharedBuffers:  "64Mb",
			HugePages: HugePages{
				Enabled: "false",
				Size:    "2Mi",
				Memory:  "",
			},
		},
		Storage: Storage{
			Enabled:         "false",
			CcpStorageClass: "",
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
				Image:            "quay.io/external_storage/nfs-client-provisioner:latest",
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
		},
	}
}

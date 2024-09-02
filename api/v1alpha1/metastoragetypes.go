package v1alpha1

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

const (
	Running  MetaStorageProvisionerStatus = "running"
	Failed   MetaStorageProvisionerStatus = "failed"
	Pending  MetaStorageProvisionerStatus = "pending"
	Deleting MetaStorageProvisionerStatus = "deleting"
)

func init() {
	SchemeBuilder.Register(&MetaStorageProvisioner{}, &MetaStorageProvisionerList{})
}

type MetaStorageProvisionerStatus string

// MetaStorageProvisioner represents the storage provisioner to be installed
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`

type MetaStorageProvisioner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StorageProvisionerSpec   `json:"spec, omitempty"`
	Status StorageProvisionerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type MetaStorageProvisionerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []MetaStorageProvisioner `json:"items"`
}

// StorageProvisionerSpec defines the desired state of MetaStorageProvisioner
type StorageProvisionerSpec struct {
	ProvisionerType `json:",inline"`
}

// StorageProvisionerStatus defines the observed state of MetaStorageProvisioner
type StorageProvisionerStatus struct {
	Status MetaStorageProvisionerStatus `json:"status"`
}

type ProvisionerType struct {
	// only one of the following should be set
	// +optional
	NFSProvisioner *NFSProvisioner `json:"NFSProvisioner,omitempty"`
}

// NFSProvisioner represents the NFS provisioner
type NFSProvisioner struct {
	NFSServer string `json:"nfsServer"`
	NFSPath   string `json:"nfsPath"`
	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`
}

// Validate ensures provisioner type is valid
func (p *ProvisionerType) Validate() error {
	count := 0
	if p.NFSProvisioner != nil {
		count++
		if p.NFSProvisioner.NFSServer == "" {
			return fmt.Errorf("nfs server cannot be empty")
		}
		if p.NFSProvisioner.NFSPath == "" {
			return fmt.Errorf("nfs path cannot be empty")
		}
	}

	if count > 1 {
		return fmt.Errorf("only one provisioner type can be set")
	}
	return nil
}

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CnvrgInfraSpec struct {
	CnvrgInfraNs string  `json:"cnvrgInfraNs"`
	Storage      Storage `json:"storage,omitempty"`
	Istio        Istio   `json:"istio,omitempty"`
}

type CnvrgInfraStatus struct {
	Status   OperatorStatus `json:"status,omitempty"`
	Message  string         `json:"message,omitempty"`
	Progress string         `json:"progress,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.message`
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status
type CnvrgInfra struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CnvrgInfraSpec   `json:"spec,omitempty"`
	Status CnvrgInfraStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type CnvrgInfraList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CnvrgInfra `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CnvrgInfra{}, &CnvrgInfraList{})
}

func DefaultCnvrgInfraSpec() CnvrgInfraSpec {
	return CnvrgInfraSpec{
		CnvrgInfraNs: "cnvrg-infra",
		Storage:      storageDefault,
		Istio:        istioDefault,
	}
}

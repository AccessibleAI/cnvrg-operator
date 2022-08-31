package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CnvrgThirdPartySpec struct {
	ImageHub string   `json:"imageHub,omitempty"`
	Gpu      Gpu      `json:"gpu,omitempty"`
	Istio    Istio    `json:"istio,omitempty"`
	Registry Registry `json:"registry,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.message`
// +kubebuilder:subresource:status

type CnvrgThirdParty struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CnvrgThirdPartySpec `json:"spec,omitempty"`
	Status Status              `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type CnvrgThirdPartyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CnvrgThirdParty `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CnvrgThirdParty{}, &CnvrgThirdPartyList{})
}

func DefaultCnvrgThirdPartySpec() CnvrgThirdPartySpec {
	return CnvrgThirdPartySpec{
		ImageHub: "docker.io/cnvrg",
		Istio:    istioDefault,
		Registry: thirdPartyRegistryDefault,
		Gpu:      gpuDefaults,
	}
}

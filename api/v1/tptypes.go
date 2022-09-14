package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CnvrgThirdPartySpec struct {
	ImageHub string   `json:"imageHub,omitempty"`
	Nvidia   Nvidia   `json:"nvidia,omitempty"`
	Habana   Habana   `json:"habana,omitempty"`
	Metagpu  Metagpu  `json:"metagpu,omitempty"`
	Istio    Istio    `json:"istio,omitempty"`
	Registry Registry `json:"registry,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.message`
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ctp

// CnvrgThirdParty represent third party components, which might be optionally deployed
// by the cnvrg operator.
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

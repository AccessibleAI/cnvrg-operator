package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CnvrgAppSpec struct {
	ClusterDomain string             `json:"clusterDomain,omitempty"`
	ControlPlan   ControlPlan        `json:"controlPlan,omitempty"`
	Networking    CnvrgAppNetworking `json:"networking,omitempty"`
	Logging       CnvrgAppLogging    `json:"logging,omitempty"`
	Monitoring    CnvrgAppMonitoring `json:"monitoring"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.controlPlan.webapp.image`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.message`
// +kubebuilder:subresource:status
type CnvrgApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CnvrgAppSpec `json:"spec,omitempty"`
	Status Status       `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
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
		ClusterDomain: "",
		ControlPlan:   controlPlanDefault,
		Logging:       cnvrgAppLoggingDefault,
		Networking:    cnvrgAppNetworkingDefault,
		Monitoring:    cnvrgAppMonitoringDefault,
	}
}

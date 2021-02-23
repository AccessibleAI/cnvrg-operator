package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OperatorStatus string

const (
	STATUS_ERROR       OperatorStatus = "ERROR"
	STATUS_RECONCILING OperatorStatus = "RECONCILING"
	STATUS_HEALTHY     OperatorStatus = "HEALTHY"
)

type CnvrgAppSpec struct {
	CnvrgNs       string      `json:"cnvrgNs,omitempty"`
	ClusterDomain string      `json:"clusterDomain,omitempty"`
	ControlPlan   ControlPlan `json:"controlPlan,omitempty"`
	Pg            Pg          `json:"pg,omitempty"`
	Storage       Storage     `json:"storage,omitempty"`
	Networking    Networking  `json:"networking,omitempty"`
	Minio         Minio       `json:"minio,omitempty"`
}

type CnvrgAppStatus struct {
	Status   OperatorStatus `json:"status,omitempty"`
	Message  string         `json:"message,omitempty"`
	Progress string         `json:"progress,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.controlPlan.webapp.image`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.message`
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status
type CnvrgApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CnvrgAppSpec   `json:"spec,omitempty"`
	Status CnvrgAppStatus `json:"status,omitempty"`
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

var DefaultSpec = CnvrgAppSpec{
	ClusterDomain: "",
	CnvrgNs:       "cnvrg",
	Pg:            pgDefault,
	Storage:       storageDefault,
	ControlPlan:   controlPlanDefault,
	Networking:    networkingDefault,
	Minio:         minioDefaults,
}

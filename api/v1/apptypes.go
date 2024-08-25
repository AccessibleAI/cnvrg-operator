package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OperatorStatus string
type IngressType string

const (
	StatusError       OperatorStatus = "ERROR"
	StatusReconciling OperatorStatus = "RECONCILING"
	StatusHealthy     OperatorStatus = "HEALTHY"
	StatusReady       OperatorStatus = "READY"
	StatusRemoving    OperatorStatus = "REMOVING"

	IstioIngress     IngressType = "istio"
	NginxIngress     IngressType = "ingress"
	OpenShiftIngress IngressType = "openshift"
	NodePortIngress  IngressType = "nodeport"
	NoneIngress      IngressType = "none"
)

type CnvrgAppSpec struct {
	ClusterDomain         string            `json:"clusterDomain,omitempty"`
	ClusterInternalDomain string            `json:"clusterInternalDomain,omitempty"`
	ImageHub              string            `json:"imageHub,omitempty"`
	Labels                map[string]string `json:"labels,omitempty"`
	Annotations           map[string]string `json:"annotations,omitempty"`
	ControlPlane          ControlPlane      `json:"controlPlane,omitempty"`
	Registry              Registry          `json:"registry,omitempty"`
	Dbs                   Dbs               `json:"dbs,omitempty"`
	Networking            Networking        `json:"networking,omitempty"`
	SSO                   SSO               `json:"sso,omitempty"`
	Tenancy               Tenancy           `json:"tenancy,omitempty"`
	PriorityClass         PriorityClass     `json:"priorityClass,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.controlPlane.image`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.message`
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=cap

// CnvrgApp represent the cnvrg.io AI/MLOps control plane stack,
// which includes frontend and backend services & persistent workloads (DBs).
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
		ClusterDomain:         "",
		ClusterInternalDomain: "cluster.local",
		ImageHub:              "docker.io/cnvrg",
		ControlPlane:          controlPlaneDefault,
		Registry: Registry{

			Name:     "cnvrg-app-registry",
			URL:      "docker.io",
			User:     "",
			Password: "",
		},
		Dbs:        appDbsDefaults,
		Networking: networkingDefault,
		SSO:        ssoDefault,
		Tenancy: Tenancy{
			Enabled: false,
			Key:     "purpose",
			Value:   "cnvrg-control-plane",
		},
		Labels:      map[string]string{"owner": "cnvrg-control-plane"},
		Annotations: nil,
		PriorityClass: PriorityClass{
			AppClassRef: "",
			JobClassRef: "",
		},
	}
}

type Status struct {
	Status               OperatorStatus  `json:"status,omitempty"`
	Message              string          `json:"message,omitempty"`
	LastFeatureFlagsHash string          `json:"lastFeatureFlagHash"`
	Progress             int             `json:"progress,omitempty"`
	StackReadiness       map[string]bool `json:"stackReadiness,omitempty"`
}

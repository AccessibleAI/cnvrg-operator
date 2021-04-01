package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CnvrgInfraSpec struct {
	ClusterDomain     string               `json:"clusterDomain,omitempty"`
	InfraNamespace    string               `json:"infraNamespace,omitempty"`
	InfraReconcilerCm string               `json:"infraReconcilerCm,omitempty"`
	Monitoring        CnvrgInfraMonitoring `json:"monitoring,omitempty"`
	Networking        CnvrgInfraNetworking `json:"networking,omitempty"`
	Logging           CnvrgInfraLogging    `json:"logging,omitempty"`
	Registry          Registry             `json:"registry,omitempty"`
	Storage           Storage              `json:"storage,omitempty"`
	Dbs               InfraDbs             `json:"dbs,omitempty"`
	SSO               SSO                  `json:"sso,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.message`
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
type CnvrgInfra struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CnvrgInfraSpec `json:"spec,omitempty"`
	Status Status         `json:"status,omitempty"`
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
	infraDefault := CnvrgInfraSpec{
		ClusterDomain:     "",
		InfraNamespace:    "cnvrg-infra",
		InfraReconcilerCm: "infra-reconciler-cm",
		SSO:               ssoDefault,
		Storage:           storageDefault,
		Networking:        cnvrgInfraNetworkingDefault,
		Monitoring:        infraMonitoringDefault,
		Logging:           cnvrgInfraLoggingDefault,
		Registry:          registryDefault,
		Dbs:               infraDbsDefaults,
	}
	return infraDefault
}

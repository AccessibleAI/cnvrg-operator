package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var InfraReconcilerCm = "infra-reconciler-cm"

type CnvrgInfraSpec struct {
	ClusterDomain  string               `json:"clusterDomain,omitempty"`
	InfraNamespace string               `json:"infraNamespace,omitempty"`
	Monitoring     CnvrgInfraMonitoring `json:"monitoring,omitempty"`
	Networking     CnvrgInfraNetworking `json:"networking,omitempty"`
	Logging        CnvrgInfraLogging    `json:"logging,omitempty"`
	Registry       Registry             `json:"registry,omitempty"`
	Storage        Storage              `json:"storage,omitempty"`
	Dbs            InfraDbs             `json:"dbs,omitempty"`
	SSO            SSO                  `json:"sso,omitempty"`
	Gpu            Gpu                  `json:"gpu,omitempty"`
	Tenancy        Tenancy              `json:"tenancy,omitempty"`
	Labels         map[string]string    `json:"labels,omitempty"`
	Annotations    map[string]string    `json:"annotations,omitempty"`
	ImageHub       string               `json:"imageHub,omitempty"`
	ConfigReloader ConfigReloader       `json:"configReloader,omitempty"`
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
		ClusterDomain:  "",
		ImageHub:       "docker.io/cnvrg",
		InfraNamespace: "cnvrg-infra",
		SSO:            ssoDefault,
		Storage:        storageDefault,
		Networking:     cnvrgInfraNetworkingDefault,
		Monitoring:     infraMonitoringDefault,
		Logging:        cnvrgInfraLoggingDefault,
		Registry:       infraRegistryDefault,
		Dbs:            infraDbsDefaults,
		Gpu:            gpuDefaults,
		Tenancy:        tenancyDefault,
		Labels:         map[string]string{"owner": "cnvrg-control-plane"},
		Annotations:    nil,
		ConfigReloader: defaultConfigReloader,
	}
	return infraDefault
}

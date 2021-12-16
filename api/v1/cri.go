package v1

// +kubebuilder:validation:Enum=docker;containerd;cri-o;""
type CriType string

const (
	CriTypeDocker     CriType = "docker"
	CriTypeContainerd CriType = "containerd"
	CriTypeCrio       CriType = "cri-o"
)

package dumper

import "github.com/spf13/viper"

type ControlPlane struct {
	Image          string
	WildcardDomain string
	Cri            string
	RegistryUser   string
	RegisterPass   string
	Ingress        string
}

func NewControlPlane() *ControlPlane {
	return &ControlPlane{
		Image:          viper.GetString("control-plane-image"),
		WildcardDomain: "",
		Cri:            "",
		RegistryUser:   "",
		RegisterPass:   "",
		Ingress:        "",
	}
}

func (p *ControlPlane) GetCliFlags() []*Param {
	return []*Param{
		{Name: "ingress", Value: "ingress", Shorthand: "s", Usage: "must be one of the: istio|ingress|openshift|nodeport"},
		{Name: "wildcard-domain", Shorthand: "d", Value: "", Usage: "the wildcard domain for cnvrg stack deployments"},
		{Name: "control-plane-image", Shorthand: "i", Value: "", Usage: "cnvrg control plane image"},
		{Name: "registry-user", Shorthand: "u", Value: "", Usage: "docker registry user"},
		{Name: "registry-password", Shorthand: "p", Value: "", Usage: "docker registry password"},
		{Name: "Cri", Shorthand: "c", Value: "containerd", Usage: "container runtime interface one of: docker|containerd|cri-o"},
	}
}

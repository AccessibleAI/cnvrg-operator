package dumper

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/AccessibleAI/cnvrg-operator/controllers"
	"github.com/AccessibleAI/cnvrg-operator/pkg/desired"
	"github.com/AccessibleAI/cnvrg-operator/pkg/networking"
	"github.com/AccessibleAI/cnvrg-operator/pkg/priorityclass"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type Networking struct {
	infra *mlopsv1.CnvrgInfra
	state []*desired.State
}

func NewNetworking(ns string) *Networking {
	infra := &mlopsv1.CnvrgInfra{Spec: mlopsv1.DefaultCnvrgInfraSpec()}
	infra.Spec.InfraNamespace = ns
	if err := controllers.CalculateAndApplyInfraDefaults(infra, &infra.Spec, nil); err != nil {
		log.Fatal(err)
	}
	return &Networking{
		infra: infra,
	}
}

func (n *Networking) BuildState() error {
	n.infra.Spec.Networking.Istio.Enabled = true
	n.infra.Spec.CnvrgAppPriorityClass = mlopsv1.PriorityClass{Name: "cnvrg-apps-infra", Value: 2000000, Description: "cnvrg control plane apps priority class"}

	n.state = networking.InfraNetworkingState(n.infra)
	n.state = append(n.state, networking.IstioCrds()...)
	n.state = append(n.state, priorityclass.State()...)
	return nil
}

func (n *Networking) Dump(preserveTemplatesDir bool) error {
	if err := os.RemoveAll(viper.GetString("templates-dump-dir")); err != nil {
		log.Fatal(err)
	}
	for _, o := range n.state {

		if o.TemplateData == nil {
			o.TemplateData = n.infra
		}
		if err := o.GenerateDeployable(); err != nil {
			return err
		}
		if err := o.DumpTemplateToFile(preserveTemplatesDir); err != nil {
			return err
		}
	}
	return nil
}
